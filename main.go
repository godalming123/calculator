package main

// TODO: Give a clear error when the calculation is invalid, EG: 5+4-

import (
	"fmt"
	"strconv"
	"strings"
)

// Logging helper
var logLevel = 0
func log(toLog ...any) {
	print("\033[90m", strings.Repeat("│ ", logLevel), "\033[0m")
	fmt.Println(toLog...)
}

// Checks if a charecter is a digit
func isDigit(c byte) bool {
	if (c >= '0' && c <= '9') || c == '.' {
		return true
	}
	return false
}

// Converts float64 into string
func float64ToString(f float64) string {
	return fmt.Sprintf("%f", f)
}

// Converts a string into a float64
func stringToFloat64(string string) (float64, error) {
	if string[0] == '+' {
		return stringToFloat64(string[1:])
	} else if string[0] == '-' {
		result, err := stringToFloat64(string[1:])
		return -result, err
	} else {
		return strconv.ParseFloat(string, 64)
	}
}

// Checks if an expression is just a number
func isOneNumber(expression string) bool {
	charecterIndex := 0
	for charecterIndex < len(expression) && (expression[charecterIndex] == '+' || expression[charecterIndex] == '-') {
		charecterIndex++
	}
	for charecterIndex < len(expression) && isDigit(expression[charecterIndex]) {
		charecterIndex++
	}
	return charecterIndex == len(expression)
}

// Finds the start of the value before the operationIndex given within the expression given, EG:
// 6+--15*-21 < expression
//	 |   ^ operationIndex
//	 ^ returned value
func findValueBeforeOperation(operationIndex int, expression string) int {
	operationIndex--
	for operationIndex >= 0 && isDigit(expression[operationIndex]) {
		operationIndex--
	}
	for operationIndex >= 0 && !isDigit(expression[operationIndex]) {
		operationIndex--
	}
	if operationIndex == -1 {
		return 0
	}
	return operationIndex + 2
}

// Finds the end of the value after the operationIndex given within the expression given, EG:
// 6+--15*-21 < expression
//	|   ^ returned value
//	^ operationIndex
func findValueAfterOperation(operationIndex int, expression string) int {
	for operationIndex < len(expression) && !isDigit(expression[operationIndex]) {
		operationIndex++
	}
	for operationIndex < len(expression) && isDigit(expression[operationIndex]) {
		operationIndex++
	}
	return operationIndex
}

// Simplifies an expression around the operation index EG:
// 6+--15*-21 < expression
//       ^ operationIndex
// 6+-315 < returned value
func simplifyExpression(expression string, operationIndex int) (string, error) {
	log("Simplifying operation " + string(expression[operationIndex]) + " at index " + strconv.Itoa(operationIndex) + " in the expression " + expression)

	// Find the start and end of the calculation
	startOfCalculation := findValueBeforeOperation(operationIndex, expression)
	endOfCalculation := findValueAfterOperation(operationIndex, expression)

	// Find the value of the start of the calculation
	valueOfExpressionBeforeOperator, err := stringToFloat64(expression[startOfCalculation:operationIndex])
	if err != nil {
		return "", err
	}

	// Find the value of the end of the calculation
	valueOfExpressionAfterOperator, err := stringToFloat64(expression[operationIndex+1:endOfCalculation])
	if err != nil {
		return "", err
	}

	// Calculate the result of the expression
	var calculationResult float64
	if expression[operationIndex] == '*' {
		calculationResult = valueOfExpressionBeforeOperator*valueOfExpressionAfterOperator
	} else if expression[operationIndex] == '/' {
		calculationResult = valueOfExpressionBeforeOperator/valueOfExpressionAfterOperator
	} else if expression[operationIndex] == '+' {
		calculationResult = valueOfExpressionBeforeOperator+valueOfExpressionAfterOperator
	} else if expression[operationIndex] == '-' {
		calculationResult = valueOfExpressionBeforeOperator-valueOfExpressionAfterOperator
	}

	// Return
	expression = expression[:startOfCalculation] + float64ToString(calculationResult) + expression[endOfCalculation:]
	log("Simplified the operation, the expression is now " + expression)
	return expression, nil
}

// Finds the simplified result of an expression that does not have brackets
func calculateExpressionWithoutBrackets(expression string) (float64, error) {
	var err error
	
	// Simplify * and /
	for index := 0; index < len(expression); index++ {
		if expression[index] == '*' || expression[index] == '/' {
			expression, err = simplifyExpression(expression, index)
			if err != nil {
				return 0, err
			}
			index = 0
		}
	}

	// Simplify + and -
	for true {
		// Consider returning
		if isOneNumber(expression) {
			goto returnCode
		}

		// Find the first occurance of + or - that isn't just the starting + or - in a expression (EG: -4+2)
		index := 0
		for expression[index] == '+' || expression[index] == '-' {
			index++
		}
		for index < len(expression) {
			if expression[index] == '+' || expression[index] == '-' {
				break
			}
			index++
		}

		// Simplify expression
		expression, err = simplifyExpression(expression, index)
		if err != nil {
			return 0, err
		}
	}

	returnCode:
	  return stringToFloat64(string(expression))
}

// Finds the simplified result of an expression that has brackets
func calculateExpression(expression string) (float64, error) {
	expression = "(" + expression + ")"
	log("Calculating the result of the expression " + expression)
    logLevel++
	for !isOneNumber(expression) {
		// Find the last open bracket
		lastOpenBracket := len(expression) - 1
		for lastOpenBracket > 0 {
			if expression[lastOpenBracket] == '(' {
				break
			}
			lastOpenBracket--
		}

		// Find the first close bracket after the last open bracket
		firstCloseBracket := lastOpenBracket
		for firstCloseBracket < len(expression) {
			if expression[firstCloseBracket] == ')' {
				break
			}
			firstCloseBracket++
		}

		// Simplify the expression inside the brackets
		log("Simplifying the brackets " + expression[lastOpenBracket:firstCloseBracket+1] + " in the expression " + expression)
		logLevel++
		expressionResult, err := calculateExpressionWithoutBrackets(expression[lastOpenBracket+1:firstCloseBracket])
		if err != nil {
			return 0, err
		}
		expression = strings.ReplaceAll(
			expression,
			expression[lastOpenBracket:firstCloseBracket+1],
			float64ToString(expressionResult),
		)
		logLevel--
		log("Simplified the brackets, the expression is now " + expression)
	}
	logLevel--
	log("Expression result calculated to be " + expression)
	return stringToFloat64(expression)
}

func main() {
	_, err := calculateExpression("(5+4)*-2-+3")
	if err != nil {
		println(err.Error())
	}
}
