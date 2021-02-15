package main

import "testing"

func TestStore(t *testing.T) {
	f := NewFrame(
		NewCode(
			Stack{
				Val: []interface{}{
					// x = 3
					OpConst,
					OpStore,
					// y = 4
					OpConst,
					OpStore,
					// z = x + y
					OpLoad,
					OpLoad,
					OpAdd,
					OpStore,
					// none
					OpConst,
					OpRet,
				},
			},
			Stack{
				Val: []interface{}{
					NoneObj{},
					NumObj{3},
					NumObj{4},
				},
			},
			Stack{
				Val: []interface{}{
					"x",
					"y",
					"z",
				},
			},
			Stack{
				Val: []interface{}{
					1,
					0,
					2,
					1,
					0,
					1,
					2,
					0,
				},
			},
		),
	)
	vm := Vm{
		frames: []*Frame{
			f,
		},
	}
	if err := vm.eval(); err != nil {
		panic(err)
	}
}

func TestPut(t *testing.T) {
	f := NewFrame(
		NewCode(
			Stack{
				Val: []interface{}{
					OpConst,
					OpPut,
					OpPop,
					OpConst,
					OpRet,
				},
			},
			Stack{
				Val: []interface{}{
					NoneObj{},
					StrObj{"Hello World!!"},
				},
			},
			Stack{
				Val: []interface{}{
				},
			},
			Stack{
				Val: []interface{}{
					1,
					0,
				},
			},
		),
	)
	vm := Vm{
		frames: []*Frame{
			f,
		},
	}
	if err := vm.eval(); err != nil {
		panic(err)
	}
}

func TestFourOperations(t *testing.T) {
	f := NewFrame(
		NewCode(
			Stack{
				Val: []interface{}{
					OpConst,
					OpConst,
					OpAdd,
					OpPut,
					OpConst,
					OpMul,
					OpPut,
					OpConst,
					OpRet,
				},
			},
			Stack{
				Val: []interface{}{
					NoneObj{},
					NumObj{34},
					NumObj{66},
					NumObj{8},
				},
			},
			Stack{
				Val: []interface{}{
				},
			},
			Stack{
				Val: []interface{}{
					1,
					2,
					3,
					0,
				},
			},
		),
	)
	vm := Vm{
		frames: []*Frame{
			f,
		},
	}
	if err := vm.eval(); err != nil {
		panic(err)
	}
}

func TestStoreAndOutput(t *testing.T) {
	f := NewFrame(
		NewCode(
			Stack{
				Val: []interface{}{
					OpConst,
					OpStore,
					OpLoad,
					OpPut,
					OpPop,
					OpConst,
					OpRet,
				},
			},
			Stack{
				Val: []interface{}{
					NoneObj{},
					StrObj{"Hello World!!"},
				},
			},
			Stack{
				Val: []interface{}{
					"x",
				},
			},
			Stack{
				Val: []interface{}{
					1,
					0,
					0,
					0,
				},
			},
		),
	)
	vm := Vm{
		frames: []*Frame{
			f,
		},
	}
	if err := vm.eval(); err != nil {
		panic(err)
	}
}

func TestFunctionCall(t *testing.T) {
	f := NewFrame(
		NewCode(
			Stack{
				Val: []interface{}{
					// x = 3
					OpConst,
					OpStore,
					// foo()
					OpConst,
					OpCall,
					OpPop,
					// foo()
					OpConst,
					OpCall,
					OpPop,
					// x
					OpLoad,
					OpPut,
					OpPop,
					// none
					OpConst,
					OpRet,
				},
			},
			Stack{
				Val: []interface{}{
					NoneObj{},
					NumObj{3},
					FuncObj{
						name: "foo",
						args: 0,
						code: NewCode(
							Stack{
								Val: []interface{}{
									// x = 5
									OpConst,
									OpStore,
									// x
									OpLoad,
									OpPut,
									OpPop,
									// Drift!!
									OpConst,
									OpPut,
									OpPop,
									// none
									OpConst,
									OpRet,
								},
							},
							Stack{
								Val: []interface{}{
									NoneObj{},
									NumObj{5},
									StrObj{"Drift!!"},
								},
							},
							Stack{
								Val: []interface{}{
									"x",
								},
							},
							Stack{
								Val: []interface{}{
									1, // CONST
									0, // STORE
									0, // LOAD
									2, // CONST
									0, // CONST
								},
							},
						),
					},
				},
			},
			Stack{
				Val: []interface{}{
					"x",
				},
			},
			Stack{
				Val: []interface{}{
					1, // CONST
					0, // STORE
					2, // CONST
					0, // CALL
					2, // CONST
					0, // CALL
					0, // LOAD
					0, // CONST
				},
			},
		),
	)
	vm := Vm{
		frames: []*Frame{
			f,
		},
	}
	if err := vm.eval(); err != nil {
		panic(err)
	}
}

func TestFunctionArgs(t *testing.T) {
	f := NewFrame(
		NewCode(
			Stack{
				Val: []interface{}{
					OpConst, // 3
					OpConst, // OK!!
					// foo
					OpConst,
					OpCall,
					OpPop,
					// foo
					OpConst, // 6.3
					OpConst, // 3.7
					// foo
					OpConst,
					OpCall,
					OpPop,
					// none
					OpConst,
					OpRet,
				},
			},
			Stack{
				Val: []interface{}{
					NoneObj{},
					NumObj{3},
					StrObj{"OK!!"},
					NumObj{6.3},
					NumObj{3.7},
					FuncObj{
						name: "foo",
						args: 2,
						code: NewCode(
							Stack{
								Val: []interface{}{
									// x
									OpLoad,
									OpPut,
									OpPop,
									// y
									OpLoad,
									OpPut,
									OpPop,
									// none
									OpConst,
									OpRet,
								},
							},
							Stack{
								Val: []interface{}{
									NoneObj{},
								},
							},
							Stack{
								Val: []interface{}{
									"x",
									"y",
								},
							},
							Stack{
								Val: []interface{}{
									0,
									1,
									0,
								},
							},
						),
					},
				},
			},
			Stack{
				Val: []interface{}{
				},
			},
			Stack{
				Val: []interface{}{
					1, // CONST
					2, // CONST
					5, // CONST
					2, // CALL
					3, // CONST
					4, // CONST
					5, // CONST
					2, // CALL
					0, // CONST
				},
			},
		),
	)
	vm := Vm{
		frames: []*Frame{
			f,
		},
	}
	if err := vm.eval(); err != nil {
		panic(err)
	}
}

func TestFunctionReturn(t *testing.T) {
	f := NewFrame(
		NewCode(
			Stack{
				Val: []interface{}{
					OpConst, // 5
					OpPut,
					OpConst, // plus
					OpCall,  // 1
					OpPut,
					OpPop,
					// none
					OpConst,
					OpRet,
				},
			},
			Stack{
				Val: []interface{}{
					NoneObj{},
					NumObj{5},
					FuncObj{
						name: "plus",
						args: 1,
						code: NewCode(
							Stack{
								Val: []interface{}{
									OpLoad, // num
									OpPut,
									OpAddO,
									OpPut,
									OpRet,
								},
							},
							Stack{
								Val: []interface{}{
									NoneObj{},
								},
							},
							Stack{
								Val: []interface{}{
									"num",
								},
							},
							Stack{
								Val: []interface{}{
									0, // LOAD
									0, // CONST
								},
							},
						),
					},
				},
			},
			Stack{
				Val: []interface{}{
				},
			},
			Stack{
				Val: []interface{}{
					1, // CONST
					2, // CONST
					1, // CALL
					0,
				},
			},
		),
	)
	vm := Vm{
		frames: []*Frame{
			f,
		},
	}
	if err := vm.eval(); err != nil {
		panic(err)
	}
}
