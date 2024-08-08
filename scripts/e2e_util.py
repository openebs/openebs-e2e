'''
Utility functions for e2e python3 scripts
'''
import json
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

def json_load(filename):
    '''
    utility function to load a json file
    '''
    with open(filename, encoding='UTF-8') as jfp:
        return json.load(jfp)


def json_save(obj, filename):
    '''
    utility function to save and object to a json file
    '''
    with open(filename, 'w', encoding='UTF-8') as ofp:
        json.dump(obj, ofp, indent=4, sort_keys=True)


def load_testplan(testplan):
    '''
    load a testplan, return testcases and array of testplan meta data
    '''
    tplan = yaml_load(e2e_root + f'/testplans/{testplan}.yaml')
    testcases = tplan.get('testcases', {})
    metas = [tplan['meta']]
    for inc_tp_name in tplan['meta'].get('include',[]):
        included_testcases, included_meta = load_testplan(inc_tp_name)
        testcases.update(included_testcases)
        metas.extend(included_meta)
    return testcases, metas


def testcases_from_testplan(testplan):
    '''
    yields a filtered set testcases found in the testplan
    filters out BeforeSuite and Basic Install Suite
    '''
    tpl, _ = load_testplan(testplan)
    for tcs in [tpl[k]['testcase'] for k in tpl]:
        if tcs.endswith('BeforeSuite') or tcs.endswith('AfterSuite'):
            continue
        yield tcs


def testcases_from_testplans(testplans):
    '''
    yields a filtered set testcases found in the testplans
    filters out BeforeSuite and Basic Install Suite
    '''
    for testplan in testplans:
        yield from testcases_from_testplan(testplan)
