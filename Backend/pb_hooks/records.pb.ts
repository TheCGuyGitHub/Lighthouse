/// <reference path="../pb_data/types.d.ts" />

onCollectionCreate((e) => {
    // In the Future we will set some permission stuff
    e.next()
}, "users");