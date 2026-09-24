package presence_test

import (
	"encoding/json"
	"fmt"

	"github.com/pivaldi/presence/v2"
)

// Example demonstrates basic usage of presence values.
func Example() {
	// Create presence values
	name := presence.FromValue("John Doe")
	age := presence.FromValue(30)
	email := presence.Null[string]() // Explicitly null

	// Check states
	fmt.Println("Name is null:", name.IsNull())
	fmt.Println("Email is null:", email.IsNull())

	// Get values
	if !age.IsNull() {
		fmt.Println("Age:", *age.GetValue())
	}

	// Output:
	// Name is null: false
	// Email is null: true
	// Age: 30
}

// ExampleFromValue demonstrates creating presence values from concrete values.
func ExampleFromValue() {
	name := presence.FromValue("Alice")
	age := presence.FromValue(25)
	active := presence.FromValue(true)

	fmt.Println(*name.GetValue())
	fmt.Println(*age.GetValue())
	fmt.Println(*active.GetValue())

	// Output:
	// Alice
	// 25
	// true
}

// ExampleNull demonstrates creating explicit null presence values.
func ExampleNull() {
	email := presence.Null[string]()
	age := presence.Null[int]()

	fmt.Println("Email is null:", email.IsNull())
	fmt.Println("Age is null:", age.IsNull())

	// Output:
	// Email is null: true
	// Age is null: true
}

// ExampleOf_json demonstrates JSON marshaling and unmarshaling.
func ExampleOf_json() {
	type User struct {
		Name  presence.Of[string] `json:"name"`
		Email presence.Of[string] `json:"email"`
		Age   presence.Of[int]    `json:"age"`
	}

	// Create user with some null fields
	user := User{
		Name:  presence.FromValue("John"),
		Email: presence.Null[string](),
		Age:   presence.FromValue(30),
	}

	// Marshal to JSON
	data, _ := json.Marshal(user)
	fmt.Println(string(data))

	// Unmarshal from JSON
	jsonStr := `{"name":"Jane","email":"jane@example.com","age":null}`
	var user2 User
	json.Unmarshal([]byte(jsonStr), &user2)

	fmt.Println(*user2.Name.GetValue())
	fmt.Println(*user2.Email.GetValue())
	fmt.Println("Age is null:", user2.Age.IsNull())

	// Output:
	// {"name":"John","email":null,"age":30}
	// Jane
	// jane@example.com
	// Age is null: true
}

// ExampleOf_threeState demonstrates the three-state model (unset, null, value).
func ExampleOf_threeState() {
	// Unset field - not touched
	var unset presence.Of[string]
	fmt.Println("Unset - IsUnset:", unset.IsUnset())
	fmt.Println("Unset - IsNull:", unset.IsNull())
	fmt.Println("Unset - IsSet:", unset.IsSet())

	// Null field - explicitly set to null
	null := presence.Null[string]()
	fmt.Println("Null - IsUnset:", null.IsUnset())
	fmt.Println("Null - IsNull:", null.IsNull())
	fmt.Println("Null - IsSet:", null.IsSet())

	// Value field - has concrete value
	value := presence.FromValue("John")
	fmt.Println("Value - IsUnset:", value.IsUnset())
	fmt.Println("Value - IsNull:", value.IsNull())
	fmt.Println("Value - IsValue:", value.IsValue())

	// Output:
	// Unset - IsUnset: true
	// Unset - IsNull: false
	// Unset - IsSet: false
	// Null - IsUnset: false
	// Null - IsNull: true
	// Null - IsSet: true
	// Value - IsUnset: false
	// Value - IsNull: false
	// Value - IsValue: true
}

// ExampleOf_patchAPI demonstrates using presence for PATCH API semantics.
func ExampleOf_patchAPI() {
	type UpdateUserRequest struct {
		Name  presence.Of[string] `json:"name,omitempty"`
		Email presence.Of[string] `json:"email,omitempty"`
		Age   presence.Of[int]    `json:"age,omitempty"`
	}

	// Simulate different PATCH requests
	requests := []string{
		`{}`,                                 // No updates
		`{"name":"John"}`,                    // Update name only
		`{"name":"John","age":null}`,         // Update name, clear age
		`{"name":null,"email":"new@ex.com"}`, // Clear name, update email
	}

	for _, jsonStr := range requests {
		var req UpdateUserRequest
		json.Unmarshal([]byte(jsonStr), &req)

		fmt.Println("Request:", jsonStr)
		if req.Name.IsSet() {
			if req.Name.IsNull() {
				fmt.Println("  Clear name")
			} else {
				fmt.Println("  Update name to:", *req.Name.GetValue())
			}
		}
		if req.Email.IsSet() {
			fmt.Println("  Update email to:", *req.Email.GetValue())
		}
		if req.Age.IsSet() && req.Age.IsNull() {
			fmt.Println("  Clear age")
		}
	}

	// Output:
	// Request: {}
	// Request: {"name":"John"}
	//   Update name to: John
	// Request: {"name":"John","age":null}
	//   Update name to: John
	//   Clear age
	// Request: {"name":null,"email":"new@ex.com"}
	//   Clear name
	//   Update email to: new@ex.com
}

// ExampleMap demonstrates transforming presence values.
func ExampleMap() {
	age := presence.FromValue(25)

	// Transform int to string
	ageStr := presence.Map(age, func(a int) string {
		return fmt.Sprintf("%d years old", a)
	})

	fmt.Println(*ageStr.GetValue())

	// Map on null value returns null
	nullAge := presence.Null[int]()
	nullAgeStr := presence.Map(nullAge, func(a int) string {
		return fmt.Sprintf("%d years old", a)
	})

	fmt.Println("Result is null:", nullAgeStr.IsNull())

	// Output:
	// 25 years old
	// Result is null: true
}

// ExampleMapOr demonstrates transforming with a default value.
func ExampleMapOr() {
	age := presence.FromValue(25)
	nullAge := presence.Null[int]()

	// Transform with default for null
	result1 := presence.MapOr(age, "unknown", func(a int) string {
		return fmt.Sprintf("%d years old", a)
	})
	result2 := presence.MapOr(nullAge, "unknown", func(a int) string {
		return fmt.Sprintf("%d years old", a)
	})

	fmt.Println(result1)
	fmt.Println(result2)

	// Output:
	// 25 years old
	// unknown
}

// ExampleFilter demonstrates filtering presence values based on a predicate.
func ExampleFilter() {
	age1 := presence.FromValue(25)
	age2 := presence.FromValue(15)

	// Keep only if >= 18
	adult1 := presence.Filter(age1, func(a int) bool {
		return a >= 18
	})
	adult2 := presence.Filter(age2, func(a int) bool {
		return a >= 18
	})

	fmt.Println("Age 25 passes:", adult1.IsValue())
	fmt.Println("Age 25 value:", *adult1.GetValue())
	fmt.Println("Age 15 passes:", adult2.IsValue())
	fmt.Println("Age 15 is null:", adult2.IsNull())

	// Output:
	// Age 25 passes: true
	// Age 25 value: 25
	// Age 15 passes: false
	// Age 15 is null: true
}

// ExampleOr demonstrates returning the first non-null value.
func ExampleOr() {
	preferred := presence.Null[string]()
	display := presence.Null[string]()
	defaultName := presence.FromValue("Guest")

	// Returns first non-null value
	name := presence.Or(preferred, display, defaultName)
	fmt.Println(*name.GetValue())

	// With a non-null preferred name
	preferred2 := presence.FromValue("John")
	name2 := presence.Or(preferred2, display, defaultName)
	fmt.Println(*name2.GetValue())

	// Output:
	// Guest
	// John
}

// ExampleOf_GetOr demonstrates getting values with defaults.
func ExampleOf_GetOr() {
	name := presence.FromValue("John")
	email := presence.Null[string]()

	fmt.Println(name.GetOr("Guest"))
	fmt.Println(email.GetOr("no-email@example.com"))

	// Output:
	// John
	// no-email@example.com
}

// ExampleOf_nested demonstrates nested presence structures.
func ExampleOf_nested() {
	type Address struct {
		Street presence.Of[string] `json:"street"`
		City   presence.Of[string] `json:"city"`
	}

	type User struct {
		Name    presence.Of[string]  `json:"name"`
		Address presence.Of[Address] `json:"address"`
	}

	user := User{
		Name: presence.FromValue("John"),
		Address: presence.FromValue(Address{
			Street: presence.FromValue("123 Main St"),
			City:   presence.FromValue("New York"),
		}),
	}

	data, _ := json.MarshalIndent(user, "", "  ")
	fmt.Println(string(data))

	// Output:
	// {
	//   "name": "John",
	//   "address": {
	//     "street": "123 Main St",
	//     "city": "New York"
	//   }
	// }
}

// ExampleFromPtr demonstrates creating presence values from pointers.
func ExampleFromPtr() {
	str := "Hello"
	var nilStr *string

	// From non-nil pointer
	value1 := presence.FromPtr(&str)
	fmt.Println("Has value:", value1.IsValue())
	fmt.Println("Value:", *value1.GetValue())

	// From nil pointer becomes null
	value2 := presence.FromPtr(nilStr)
	fmt.Println("Is null:", value2.IsNull())

	// Output:
	// Has value: true
	// Value: Hello
	// Is null: true
}

// ExampleFromBool demonstrates creating presence values based on a boolean condition.
func ExampleFromBool() {
	// Like map lookup: value, ok
	value, ok := map[string]string{"key": "value"}["key"]
	present1 := presence.FromBool(value, ok)

	value2, ok2 := map[string]string{}["missing"]
	present2 := presence.FromBool(value2, ok2)

	fmt.Println("Found key:", present1.IsValue())
	fmt.Println("Value:", *present1.GetValue())
	fmt.Println("Missing key is null:", present2.IsNull())

	// Output:
	// Found key: true
	// Value: value
	// Missing key is null: true
}

// ExampleSetClause demonstrates building a partial UPDATE from a PATCH request.
func ExampleSetClause() {
	type UpdateUserRequest struct {
		Name  presence.Of[string] `json:"name,omitzero"`
		Email presence.Of[string] `json:"email,omitzero"`
		Age   presence.Of[int]    `json:"age,omitzero"`
	}

	var req UpdateUserRequest
	_ = json.Unmarshal([]byte(`{"name": "John", "email": null}`), &req)

	clause, args := presence.SetClause(presence.Dollar,
		presence.Set("name", req.Name),   // set → updated
		presence.Set("email", req.Email), // null → cleared
		presence.Set("age", req.Age),     // unset → untouched
	)

	if clause == "" {
		fmt.Println("nothing to update")

		return
	}

	userID := 42
	args = append(args, userID)
	query := "UPDATE users SET " + clause + " WHERE id = " + presence.Dollar(len(args))

	fmt.Println(query)
	fmt.Println(len(args), "args")
	// db.Exec(query, args...)

	// Output:
	// UPDATE users SET name = $1, email = $2 WHERE id = $3
	// 3 args
}
