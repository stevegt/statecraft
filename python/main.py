import sys
import traceback
from sc import Machine, SCErr

VERSION = "v0.6.0"
USAGE = f'''
statecraft {VERSION}
usage: %(progname)s {{infn}} {{outfn}}
'''

def load_and_save(infn, outfn):
    try:
        with open(infn) as infh:
            m = Machine(infn)
            m.load(infh)

        suffix = outfn.split('.')[-1]

        if suffix == "dot":
            buf = m.to_dot()
        elif suffix == "go":
            buf = m.to_go()
            
        with open(outfn, 'w') as outfh:
            outfh.write(buf)
    except SCErr as e:
        print(e, file=sys.stderr)
        sys.exit(e.rc)
    except Exception as e:
        traceback.print_exc()
        sys.exit(-1)

def main():
    if len(sys.argv) != 3:
        print(USAGE % {'progname': sys.argv[0]}, file=sys.stderr)
        sys.exit(1)
    
    infn = sys.argv[1]
    outfn = sys.argv[2]
    
    load_and_save(infn, outfn)

if __name__ == "__main__":
    main()

