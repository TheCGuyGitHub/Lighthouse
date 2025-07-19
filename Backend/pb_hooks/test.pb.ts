/// <reference path="../pb_data/types.d.ts" />

routerAdd("GET", "/hello/{name}", (e) => {
    let name = e.request?.pathValue("name")
    return e.json(200, { message: "Hello " + name })
})

// After a user record is updated
onRecordAfterUpdateSuccess((e) => {
    console.log("user updated...", e.record?.get("email"))
    e.next() // allow the request to continue
}, "users")

// Runs when PocketBase starts
onBootstrap((e) => {
    e.next()
    console.log("App initialized!")
})