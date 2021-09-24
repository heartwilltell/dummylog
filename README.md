# Dummylog

Dummylog - prints silly messages to stdout or file in text of JSON format.

## ⚠️ This is not a logging library!

Main purpose is to help you test readers in your own programs.

## Usage

Write text logs to stdout:

```go
if err := dummylog.New().Start(ctx); err != nil {
return err
}
```

Write JSON structured logs to stdout:

```go
if err := dummylog.New(dummylog.WithFormat(dummylog.JSON)).Start(ctx); err != nil {
return err
}
```

Write JSON structured logs to file:

```go
file, createErr := os.Create(filePath)
if createErr != nil {
return createErr
}

dummy := dummylog.New(dummylog.WithWriter(file), dummylog.WithFormat(dummylog.JSON))
if err := dummy.Start(ctx); err != nil {
return err
}
```

## ⌨️ Command line interface

```
./dummylog run

flags: 
  -format - sets log format: 'text' or 'json'.
  -file   - sets path to file where logs will be written.
```

```
./dummylog serve

flags: 
  -file - sets path to file where logs will be written.
```

