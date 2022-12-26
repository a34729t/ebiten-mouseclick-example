# Eibten MouseClick Example

While building an async input manager that sticks input events on a channel, I found that mouse press detection behaved strangely, but only when running in a goroutine, for `inpututil.IsMouseButtonJustPressed()` which is listed as concurrent-safe in the docs. When running in main thread in `Update()` instead of in a goroutine, there is no problem!

Specifically, the async event manager produces hundreds of `Mouse just pressed` messages

    2022/12/26 12:08:09 Mouse just pressed
    ....
    2022/12/26 12:08:09 Mouse just pressed
    2022/12/26 12:08:09 Mouse just pressed
    2022/12/26 12:08:09 Mouse released at <544, 329> with 83
    2022/12/26 12:08:09 Mouse click!

Vs for the main thread version, there will only be a single `Mouse just pressed` message.

**To reproduce**, run the program as is. To see the main thread version, set the flag `USE_INPUT_MANAGER` on `main.go:11` to `false`.