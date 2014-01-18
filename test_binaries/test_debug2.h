#ifndef __TEST_DEBUG2_H_
#define __TEST_DEBUG2_H_

#include <string>

class A {
 public:
  A() : foo("A") { }
 
  virtual ~A();
  virtual std::string name() const;

 protected:
  std::string foo;
};

#endif
