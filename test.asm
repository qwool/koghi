
format ELF64 executable

macro write arg1, arg2 ;where arg1 - text and arg2 - length
{
  mov rax, 1
  mov rdi, 1
  mov rsi, arg1
  mov rdx, arg2
  syscall
}

segment readable executable
entry main
main:
  write msg, msg_len

  mov rax, 60
  mov rdi, 0
  syscall

segment readable writable
  msg db "asdfasdf", 10
  msg_len = $ - msg
