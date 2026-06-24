#!/usr/bin/env python3

import json
import sys

def main(json_files):
    #print(json_files)
    for i in json_files:
        #print("file",i)
        doc=open(i).read()
        #print (doc)
        m=json.loads(doc)
        find_types(m)
 
def find_types(m):

    if type(m) is dict:
        for k in m:
            if k=='@type':
                if type(m[k]) is list:
                    for l in m[k]:
                        print(get_prefix(l))
                else:
                    print(get_prefix(m[k]))
            else:
                find_types(m[k])
    elif type(m) is list:
        for l in m:
            find_types(l)

def get_prefix(iri):
    #print('iri',iri)
    hash_pos = iri.rfind('#')
    if hash_pos >= 0:
        return iri[:hash_pos]
    slash_pos = iri.rfind('/')
    if slash_pos < 0:
        return iri
 
    return iri[:slash_pos]


if __name__ == '__main__':
    main(sys.argv[1:])
