'''
Utility functions for e2e python3 scripts
'''

import os
import sys
import yaml

e2e_root = os.path.realpath(sys.path[0] + '/..')

def yaml_load(filename):
    '''
    utility function to load a yaml file
    '''
    with open(filename, encoding='UTF-8') as jfp:
        return yaml.safe_load(jfp.read())


def load_testplan(testplan):
    '''
    load a testplan, return testcases and array of testplan meta data
    '''
    tplan = yaml_load(e2e_root + f'/testplans/{testplan}.yaml')
    testsuites = tplan.get('testsuites', [])
    
    return testsuites

