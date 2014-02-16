package subscribe

import (
	"fmt"
	"os"
	"testing"
)

func TestYYets(t *testing.T) {
	f, err := os.Open("yyets.html")
	if err != nil {
		t.Error(err)
	}

	_, tasks, err := parse(f)
	if err != nil {
		t.Error(err)
	}

	if len(tasks) != 3 {
		t.Errorf("Expect 3 tasks but %d", len(tasks))
	}

	fmt.Printf("%v\n", tasks[0])
	fmt.Printf("%v\n", tasks[1])
	fmt.Printf("%v\n", tasks[2])
}

func TestYYets2(t *testing.T) {
	f, err := os.Open("b.html")
	if err != nil {
		t.Error(err)
	}

	_, tasks, err := parse(f)
	if err != nil {
		t.Error(err)
	}

	if len(tasks) != 4 {
		t.Errorf("Expect 4 tasks but %d", len(tasks))
	}

	fmt.Printf("%v\n", tasks[0])
	fmt.Printf("%v\n", tasks[1])
	fmt.Printf("%v\n", tasks[2])
	fmt.Printf("%v\n", tasks[3])
}
