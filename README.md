# 🤪 Dummylog

Dummylog - prints silly messages to stdout or file in text or JSON format.
___

## ⚠️ This is not a logging library!

Main purpose is to help you test readers in your own programs.
___

## Usage

Write text logs to stdout:

```
if err := dummylog.New().Start(ctx); err != nil {
    return err
}
```

Write JSON structured logs to stdout:

```
if err := dummylog.New(dummylog.WithFormat(dummylog.JSON)).Start(ctx); err != nil {
    return err
}
```

Write JSON structured logs to file:

```
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

Run dummylog as a process to prints messages to file or to stdout:

```
./dummylog run

flags: 
  -format - sets log format: 'text' or 'json'. Default is text.
  -file   - sets path to file where logs will be written, if
            flag is not present, logs will written to stdout.
```

Run dummylog as a server which listen to POST :8080/say and prints request body as message to file or to stdout:

```
./dummylog serve

flags: 
    -serveraddr - sets HTTP server address. Default is ':8080'.
    -file       - sets path to file where logs will be written, if 
                  flag is not present, logs will written to stdout.
```

