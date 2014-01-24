int foo1(void) {
  return 1;
}

extern int foo2(void);

int foo3(void) {
  return 3;
}

extern int foo4(void);

int foo5(void) {
  return 5;
}

extern int foo6(void);

int main(int argc, char* argv[]) {
  switch(argc) {
    case 1: return foo1();
    case 2: return foo2();
    case 3: return foo3();
    case 4: return foo4();
    case 5: return foo5();
    case 6: return foo6();
    default: return -1;
  }
}
