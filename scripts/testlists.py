#!/usr/bin/env python3
"""
Load a test list definition yaml file and
output  the list of tests for a specified profile
The default delimiter is a space
Optionally
    - change the delimiter
    - sort the list of tests alphbetically
    - return a list where each test is bracketed by install and uninstall
    - return a list which starts with install and ends with uninstall
"""

import e2e_util

def main(args):
    """
    The real main function
    """
    tests = e2e_util.load_testplan(args.testplan)
    if args.install:
        tests.insert(0, 'install')
    if args.uninstall:
        tests.append('uninstall')
    for test in tests:
    # Code to execute for each element
        print(test)


if __name__ == '__main__':
    from argparse import ArgumentParser

    parser = ArgumentParser()
    parser.add_argument('--testplan', dest='testplan', default=None, required=True,
                        help='testplan')
    parser.add_argument('--install', dest='install', action='store_true',
                        default=None,
                        help='Add install before test list')
    parser.add_argument('--uninstall', dest='uninstall', action='store_true',
                        default=None,
                        help='Add uninstall after test list')

    main(parser.parse_args())
