package module

import (
	"context"
	"fmt"
	"io"
	"log"
	"reflect"

	"golang.org/x/sync/errgroup"
)

type Runner interface {
	Run(ctx context.Context) (err error)
}

type Initializer interface {
	Init(ctx context.Context) (err error)
}

type PreInitializer interface {
	PreInit(ctx context.Context, field any) (err error)
}

type PostInitializer interface {
	PostInit(ctx context.Context, field any) (err error)
}

func Run[App any](ctx context.Context) (err error) {
	var app App
	if err = doInit(ctx, &app); err != nil {
		return
	}
	return doRun(ctx, &app)
}

func doInit(ctx context.Context, app any) (err error) {
	for _, field := range getFields(app) {
		if err = tryPreInit(ctx, app, field); err != nil {
			return
		}
		if err = tryInit(ctx, field); err != nil {
			return fmt.Errorf("init error at field %q: %w", field.Name, err)
		}
		if err = tryPostInit(ctx, app, field); err != nil {
			return
		}
	}
	return
}

func doRun(ctx context.Context, app any) (err error) {
	group, ctx := errgroup.WithContext(ctx)
	for _, field := range getFields(app) {
		field := field
		defer func() { tryClose(field) }()
		group.Go(func() (err error) { return tryRun(ctx, field) })
	}
	return group.Wait()
}

func requiredAppPointer(app any) (_ error) {
	if reflect.TypeOf(app).Kind() == reflect.Pointer {
		return
	}
	return fmt.Errorf("app must be a pointer")
}

func tryPreInit(ctx context.Context, app any, field fieldInfo) (err error) {
	if pi, ok := app.(PreInitializer); ok {
		if err = pi.PreInit(ctx, field.Value); err != nil {
			return
		}
	}
	return
}

func tryPostInit(ctx context.Context, app any, field fieldInfo) (err error) {
	if pi, ok := app.(PostInitializer); ok {
		if err = pi.PostInit(ctx, field.Value); err != nil {
			return
		}
	}
	return
}

func tryInit(ctx context.Context, field fieldInfo) (err error) {
	if m, ok := field.Value.(Initializer); ok {
		log.Printf("[module] Init: %s (%T)", field.Name, field.Value)
		return m.Init(ctx)
	}
	return
}

func tryRun(ctx context.Context, field fieldInfo) (err error) {
	if m, ok := field.Value.(Runner); ok {
		log.Printf("[module] Run: %s (%T)", field.Name, field.Value)
		return m.Run(ctx)
	}
	return
}

func tryClose(field fieldInfo) {
	if m, ok := field.Value.(io.Closer); ok {
		log.Printf("[module] Close: %s (%T)", field.Name, field.Value)
		if err := m.Close(); err != nil {
			log.Printf("[module] Close error: %s (%T)", field.Name, field.Value)
		}
	}
}

type fieldInfo struct {
	Index int
	Name  string
	Value any
	Tag   reflect.StructTag
}

func getFields(v any) (fields []fieldInfo) {
	rv := reflect.ValueOf(v).Elem()
	num := rv.NumField()
	fields = make([]fieldInfo, num)

	for i := 0; i < num; i++ {
		rType := rv.Type().Field(i)
		field := rv.Field(i)
		if field.CanAddr() {
			field = field.Addr()
		}
		fields[i] = fieldInfo{
			Index: i,
			Name:  rType.Name,
			Tag:   rType.Tag,
			Value: field.Interface(),
		}
	}
	return
}
