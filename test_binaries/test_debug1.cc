
#include <iostream>
#include <setjmp.h>
#include <string>

#include "test_debug2.h"
#include "test_debug3.h"

class B : public A {
 public:
  B() {
    extension = ":B";
  }
  virtual ~B();

  virtual std::string name() const;

 private:
  std::string extension;
};

B::~B() {
}

std::string B::name() const {
  return foo + extension;
}

int main(int argc, char* argv[]) {
  jmp_buf jb;
  setjmp(jb);
  A a;
  B b;
  C c;
  std::cout << "A: " << a.name() << "\n";
  std::cout << "B: " << b.name() << "\n";
  std::cout << "C: " << c.name() << "\n";
  return 0;
}
