import {FileWithHandle} from "browser-fs-access";

export class WrapFileWithHandle {
    handle: FileWithHandle;

    constructor(handle: FileWithHandle) {
        this.handle = handle;
    }

    toString() {
        if (this.handle) {
            return this.handle.name;
        } else {
            return undefined;
        }
    }
}
