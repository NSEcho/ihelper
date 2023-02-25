rpc.exports = {
    download_file(targetDir, filename) {
        var NSDocumentDirectory = 9;
        var NSLibraryDirectory = 5;

        var addr = Module.findExportByName(null, "NSSearchPathForDirectoriesInDomains");
        var NSSearchPathForDirectoriesInDomains = new NativeFunction(addr, 'pointer', ['int', 'int', 'int']);

        var dir = "";
        switch (targetDir) {
            case "B":
                var bd = ObjC.classes.NSBundle.mainBundle().bundleURL().toString().slice(7);
                bd += filename;
                dir = bd;
                break;
            case "D":
                var dirs = ObjC.Object(NSSearchPathForDirectoriesInDomains(NSDocumentDirectory, 1, 1));
                dir = dirs.objectAtIndex_(0).toString();
                break;
            case "L":
                var dirs = ObjC.Object(NSSearchPathForDirectoriesInDomains(NSLibraryDirectory, 1, 1));
                dir = dirs.objectAtIndex_(0).toString();
                break;
            default:
                console.log("unsupported directory flag");
                return;
        }

        var path = dir + "/" + filename;
        var dt = ObjC.classes.NSData.alloc().initWithContentsOfFile_(path);
        var arr = Memory.readByteArray(dt.bytes(), dt.length());
        send(filename, arr);
    },
    download_bin() {
        var execPath = ObjC.classes.NSBundle.mainBundle().executablePath();
        var dt = ObjC.classes.NSData.alloc().initWithContentsOfFile_(execPath);
        var arr = Memory.readByteArray(dt.bytes(), dt.length());
        send(execPath.toString(), arr);
    },
}