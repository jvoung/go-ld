#include <math.h>
#include <setjmp.h>
#include <stdio.h>
#include <stdlib.h>

typedef struct fp {
  int (*foo)(int);
  int *p;
} fp;

int x = 99;
int y = 0;
int Func(int a) { if (a) return x; else return y; }
int Func2(int b) { if (b) return y; else return x; }
int (*g_fp)(int) = &Func2;
int* const g_p = &y;

int __attribute__((noinline))Bar(fp* z, int p) {
  if (z->foo == &Func) {
    puts("Func!\n");
  } else {
    puts("Not Func!\n");
  }
  if (z->p)
    return z->foo(*z->p);
  else
    return z->foo(p);
}

void __attribute__((noinline)) Baz(fp* z, int p) {
  switch(p) {
    case 1: z->foo = &Func;  z->p = NULL; puts("Set Func!\n"); break;
    case 2: z->foo = &Func2; z->p = &x;   puts("Not Func!\n"); break;
    case 3: z->foo = &Func2; z->p = &y;   puts("3 it is!\n"); break;
    case 4: z->foo = g_fp;   z->p = g_p;  puts("4 it is!\n"); break;
    default: exit(1);
  }
}

void __attribute__((noinline)) Baz2(fp* z, int p) {
  switch(p) {
    case 1: z->foo = &Func;  z->p = NULL; puts("Set Func!\n"); break;
    case 2: z->foo = &Func2; z->p = &x;   puts("Not Func!\n"); break;
    case 3: z->foo = &Func2; z->p = &y;   puts("3 it is!\n"); break;
    case 4: z->foo = &Func;  z->p = NULL; puts("4 it is!\n"); break;
    default: exit(1);
  }
}

int main(int argc, char* argv[]) {
  jmp_buf jb;
  fp my_fp = { 0, 0 };
  setjmp(jb);
  Baz(&my_fp, argc);
  Baz2(&my_fp, argc);
  return Func(argc) || Bar(&my_fp, (int)((argc + 10.0)/(M_PI)));
}
