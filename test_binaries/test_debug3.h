#ifndef __TEST_DEBUG3_H_
#define __TEST_DEBUG3_H_

#include <string>

#include "test_debug2.h"

class C : public A {
 public:
  C() {
    foo = "C";
  }
 
  virtual ~C();
};

#endif
