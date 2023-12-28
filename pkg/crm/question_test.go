package crm

import (
	"testing"
)

func TestQuestionFunction(t *testing.T) {
	// Aquí es donde inicializarías cualquier cosa que tu función necesite.

	// Llama a la función que estás probando y guarda el resultado
	result := QuestionFunction()

	// Comprueba si el resultado es lo que esperas.
	if result != expected {
		t.Errorf("QuestionFunction() = %v; want %v", result, expected)
	}
}