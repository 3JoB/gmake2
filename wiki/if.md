# IF

<strong>Example</strong>

```yml
# Equal to
all: |
  @if windows == windows then iaw

iaw: |
  @echo i am windows!!!


# or
all: |
  @if linux == windows then iaw or ial

iaw: |
  @echo i am windows!!!

ial: |
  @echo i am linux!!!


# Not equal to
all: |
  @if windows2 != windows then iaw

iaw: |
  @echo i am not windows!!!


# Greater than
all: |
  @if 2 > 1 then iaw

iaw: |
  @echo i am 2!!!

# Greater than or equal
all: |
  @if 2 >= 2 then iaw

iaw: |
  @echo i am 2!!!

# Smaller than
all: |
  @if 1 < 2 then iaw

iaw: |
  @echo i am 1!!!

# Less than or equal to
all: |
  @if 1 <= 1 then iaw

iaw: |
  @echo i am 1!!!
```