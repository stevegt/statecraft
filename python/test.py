import os
import unittest
from sc import Machine, SCErr

from difflib import unified_diff

REGEN = False

class TestMachine(unittest.TestCase):

    def load_and_convert(self, infn, reffn, convert_func):
        with open(infn) as infh:
            m = Machine(infn)
            m.load(infh)
   
            got = convert_func(m)

        if REGEN:
            with open(reffn, 'w') as outfh:
                outfh.write(got)

        with open(reffn) as infh:
            ref = infh.read()

        self.assertMultiLineEqual(ref, got, '\n' + '\n'.join(
            unified_diff(ref.splitlines(), got.splitlines(), fromfile='expected', tofile='actual')))

    def test_dot(self):
        self.load_and_convert("example/stoplight/car/car.statecraft", "testdata/car.dot", lambda m: m.to_dot())

    def test_go(self):
        self.load_and_convert("example/stoplight/car/car.statecraft", "testdata/car.go", lambda m: m.to_go())

if __name__ == '__main__':
    unittest.main()
