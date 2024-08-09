#!/usr/bin/env python3
"""
exec-tests.sh pre-processor to handle upgrade tests
if --installtag is set invokes e2e-test.sh twice
 1) with --tag <installtag-value> --tests install
 2) with --tag <tag value> --tests <set of tests without install>
"""
import os
import re
import subprocess
import sys
import e2e_util
import testlists

if __name__ == '__main__':
    from argparse import ArgumentParser

    parser = ArgumentParser()
    parser.add_argument('--installtag', dest='installtag', default=None,
                        help='install tag')
    parser.add_argument('--tests', dest='tests', default=None,
                        help='tests')
    parser.add_argument('--tag', dest='tag', default=None,
                        help='tag')
    parser.add_argument('--profile', dest='profile', default=None,
                        help='profile')
    parser.add_argument('--testplan', dest='testplan', default=None,
                        help='testplan')

    rootpath = e2e_util.e2e_root
    args, shellscript_commandline = parser.parse_known_args()

    os.chdir(os.path.dirname(sys.path[0]))

    def quoted(cmdlinearg):
        '''
        quote parameter command line arguments to handle spaces correctly
        '''
        if cmdlinearg.startswith('-'):
            return cmdlinearg
        return f'"{cmdlinearg}"'
    shellscript_commandline = [quoted(arg) for arg in shellscript_commandline]

    shellscript_commandline.insert(0, f'{rootpath}/scripts/exec-tests.sh')
    tests = []
    if args.tests and args.profile:
        raise RuntimeError(
            'incompatible options specified: --profile & --tests')
    if args.tests:
        tests = re.split(r'[\s]+', args.tests)
    if args.profile:
        tests = testlists.get_test_list(args.profile)
        tests.insert(0, 'install')
        tests.append('uninstall')
        # for selfci runs use the selfci test configuration
        if args.profile in ['selfci', 'self_ci']:
            shellscript_commandline.extend(['--config', 'selfci_config.yaml'])

    if args.testplan in ['selfci']:
        shellscript_commandline.extend(['--config', 'selfci_config.yaml'])

    testlists = e2e_util.yaml_load(f'{rootpath}/configurations/testlists.yaml')
    upgrade_tests = list(testlists['metadata']
                         ['install_tag_override']['tests'])

    # check if upgrade tests have been specified
    upgrade = len([x for x in upgrade_tests if x in tests]) != 0
    exec_env = os.environ


    installtag = args.installtag
    if upgrade:
        # for upgrade tests do not re-schedule loki stateful set to the control node.
        exec_env['loki_on_control_node'] = 'false'

        if installtag is None:
            if args.testplan is None:
                print('unable to run upgrade tests, neither installtag or testplan is available')
                sys.exit(1)

            # retrieve the install tag from the testplan definition
            _, tp_metas = e2e_util.load_testplan(args.testplan)
            for tp_meta in tp_metas:
                try:
                    ufv = tp_meta['upgrade-from-version']
                    if installtag is not None and installtag != ufv:
                        raise RuntimeError(
                                'multiple incompatible versions for upgrade-from-versions'
                                f' found {installtag} + {ufv}')
                    installtag = ufv
                except KeyError:
                    pass

            if installtag is None:
                print('FAILURE: Cannot run upgrade test without known from version',
                      file=sys.stderr)
                sys.exit(1)

    if (upgrade or installtag is not None ) and 'install' in tests:
        # Remove the install test from the list of tests to be run
        # subsequent to this
        tests = [t for t in tests if t != 'install']
        install_test_commandline = shellscript_commandline[:]
        install_test_commandline.extend(
            ['--tag', installtag, '--tests', 'install'])
        cmdline = f"""nix-shell --run '{' '.join(install_test_commandline)}'"""
        print('About to execute:\n',cmdline)
        sys.stdout.flush()
        with subprocess.Popen(cmdline,
                              stdout=sys.stdout, shell=True, env=exec_env) as proc:
            proc.communicate()
            if proc.returncode != 0:
                sys.exit(proc.returncode)

    shellscript_commandline.extend(['--tag', args.tag])
    if tests:
        test_list = " ".join(tests)
        shellscript_commandline.extend(['--tests', f'"{test_list}"'])
    cmdline = f"""nix-shell --run '{' '.join(shellscript_commandline)}'"""
    print('About to execute:\n', cmdline)
    sys.stdout.flush()
    with subprocess.Popen(cmdline,
                          stdout=sys.stdout, shell=True, env=exec_env) as proc:
        proc.communicate()
        sys.exit(proc.returncode)
