#!/usr/bin/env python3
"""
Load a test list definition yaml file and
output  the list of tests for a specified profile
The default delimiter is a space
Optionally
    - change the delimiter
    - sort the list of tests alphbetically
    - sort the list of tests in reverse order of execution time
    - return a list where each test is bracketed by install and uninstall
    - return a list which starts with install and ends with uninstall
"""

import os
import json
import yaml

import analyse
import e2e_util

DEFAULT_LIST_DEF = f'{e2e_util.e2e_root}/configurations/testlists.yaml'


def testplan_2_testlist(testplan, e2e_repo):
    '''
    Read a testplan convert the testcases to list tests by searching
    the source tree.
    '''
    # for e2e a test is identified by the name of directory containing
    # code defining the testcases
    testcases2source = {
        k: v.split('/')[-2] for k, v in
        analyse.scrapeSources([
            os.path.realpath(e2e_repo + '/src'),
            os.path.realpath(e2e_repo + '/3rdparty'),
        ]).items()
    }

    # use a set to return a list without duplicate items
    return list(set([testcases2source[k] for k in
                     e2e_util.testcases_from_testplan(testplan)])
                )


def load_test_lists(filepath, e2e_repo):
    '''
    load the test lists, populate macro lists and add ALL macro list
    '''
    def macro_profile(profiles):
        tmp = set()
        for profile in profiles:
            if list_defs['testprofiles'][profile] is not None:
                tmp |= set(list_defs['testprofiles'][profile])
        return list(tmp)

    def grouping_fix(testprofile):
        '''
        Enumerate the test groupings and adjust the set of tests so that
        the individual tests which are part of a grouping are replaced
        by a sequence of tests (',' separated)
        '''
        if testprofile is not None:
            groupings = list_defs['metadata']['groupings']
            for group in groupings:
                for primary_test, test_seq in group.items():
                    if primary_test in testprofile:
                        grouped_tests = [primary_test]
                        grouped_tests.extend(test_seq)
                        testprofile = [
                            test for test in testprofile if test not in grouped_tests]
                        testprofile.append(','.join(grouped_tests))
        return testprofile

    with open(filepath, mode='r', encoding='UTF-8') as f_p:
        contents = f_p.read()
        list_defs = yaml.safe_load(contents)

        testplans = [f.rsplit('.', maxsplit=1)[0] for f in os.listdir(os.path.join(e2e_repo, 'testplans'))
                     if f.endswith('.yaml') and f not in [
            'deprecated.yaml',
            'common.yaml',
            'unscheduled.yaml',
            "upgrade.yaml"
        ]
        ]

        for testplan in testplans:
            if testplan not in list_defs['testprofiles']:
                list_defs['testprofiles'][testplan] = testplan_2_testlist(
                    testplan, e2e_repo)
            else:
                print(
                    f'Ignoring testplan definition {testplan} defined explicitly?')

        for k, val in list_defs['macro-profiles'].items():
            list_defs['testprofiles'][k] = macro_profile(val)

        for k, testprofile in list_defs['testprofiles'].items():
            list_defs['testprofiles'][k] = grouping_fix(testprofile)

        exclude_from_regression = []
        for entry in list_defs['metadata']["exclude_from_regression"]:
            if entry in list_defs['testprofiles']:
                exclude_from_regression.extend(
                    list_defs['testprofiles'][entry])
            else:
                exclude_from_regression.append(entry)
        list_defs['metadata']["exclude_from_regression"] = exclude_from_regression

        list_defs['testprofiles']['ALL'] = macro_profile(sorted(
            list_defs['testprofiles'].keys()))

        return list_defs


def get_test_list(profile, listdef=DEFAULT_LIST_DEF, sort_duration=False):
    """
    retrieve the list of tests defined for a profile
    optionally sorted by duration
    """
    list_defs = load_test_lists(listdef, e2e_util.e2e_root)
    durations = list_defs['metadata']['recorded_durations']
    tests = list_defs['testprofiles'][profile]

    tests = [tst for tst in sorted(tests) if tst not in [
        'install', 'uninstall']]
    if sort_duration:
        tests = sorted(tests)
        tests = sorted(
            tests, key=lambda x: durations.get(x, 0), reverse=True)
    return tests


def main(args):
    """
    The real main function
    """
    if args.install_uninstall:
        if args.install or args.uninstall:
            raise RuntimeError(
                'Incompatible options --iu with --install or --uninstall')

    tests = get_test_list(args.profile, args.listdef, args.sort_duration)

    if args.install_uninstall:
        tests = [f'install,{tst},uninstall' for tst in tests]
    if args.install:
        tests.insert(0, 'install')
    if args.uninstall:
        tests.append('uninstall')

    if args.format == 'json':
        output = json.dumps(tests, sort_keys=True, indent=4)
    elif args.format == 'yaml':
        output = yaml.dump(tests, sort_keys=True)
    else:
        output = args.separator.join(tests)

    if args.outputfile is None:
        print(output)
    else:
        with open(args.outputfile, "w", encoding="UTF-8") as fout:
            print(output, file=fout)


if __name__ == '__main__':
    from argparse import ArgumentParser

    parser = ArgumentParser()
    parser.add_argument('--profile', dest='profile', default=None, required=True,
                        help='profile')
    parser.add_argument('--lists', dest='listdef',
                        default=DEFAULT_LIST_DEF,
                        help='list definitions')
    parser.add_argument('--e2e_repo', dest='e2e_repo',
                        default=e2e_util.e2e_root,
                        help='path to mayastor-e2e repo checkout')

    parser.add_argument('--install', dest='install', action='store_true',
                        default=None,
                        help='Add install before all tests')
    parser.add_argument('--uninstall', dest='uninstall', action='store_true',
                        default=None,
                        help='Add uninstall before all tests')
    parser.add_argument('--iu', dest='install_uninstall', action='store_true',
                        default=None,
                        help='Each test is preceded with install and followed with an uninstall')

    parser.add_argument('--separator', dest='separator', default=' ',
                        help='seperator to use for plain text output')
    parser.add_argument('--sort_duration', dest='sort_duration', action='store_true',
                        default=False,
                        help='sort the list in order of execution time high to low')
    parser.add_argument('--outputfile', dest='outputfile', default=None,
                        help='outputfile')
    parser.add_argument('-o', dest='format', choices=[None, 'json', 'yaml'],
                        default=None,
                        help='formats')

    main(parser.parse_args())
