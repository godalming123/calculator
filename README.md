# Calculator
A calculator written in go that uses BODMAS to determine order of operations. All of the logic to calculate results has been design by me without help from any other people's ideas.

## Things that are supported
- Brackets with several layers of nesting
- The `+`, `-`, `*` and `/` operations
- Use the correct order of operations with BODMAS
- Multiple pluses and minuses in a row, for example: `4+-+--2 // Equals -8`

## Things that are not supported
- The `^` and `%` operations
- Proper handling for invalid expressions, for example: `5+4-`
- Exact calculations (currently approximations are made)
- Numbers that contain charecters other then `0`-`9` and `.` such as `π`