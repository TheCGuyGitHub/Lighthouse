/// <reference path="../pb_data/types.d.ts" />

onRecordCreate((e) => {
    // In the Future we will set some permission stuff
    e.next()
}, "users");