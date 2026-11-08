const std = @import("std");
const httpz = @import("httpz");
const repository = @import("repository");

pub const Server = @This();

const GetRequest = struct {
    id: i32,
};
const PutRequest = struct {
    name: []const u8,
};
const Response = struct {
    id: i32,
    name: []const u8,
};

var db: repository.Repository = undefined;

pub fn start(io: std.Io, allocator: std.mem.Allocator) !void {
    db = try repository.init();
    defer db.deinit();

    var server = try httpz.Server(void).init(
        io,
        allocator,
        .{ .address = .localhost(16000) },
        {},
    );
    defer {
        server.stop();
        server.deinit();
    }
    var router = try server.router(.{});
    router.get("/event/:id", handleGetEvent, .{});
    router.put("/event", handlePutEvent, .{});

    try server.listen();
}

fn handleGetEvent(req: *httpz.Request, res: *httpz.Response) !void {
    res.status = 200;
    const request = (try req.json(GetRequest)).?;
    std.debug.print("got GET request: {any}\n", .{request});
    const event = db.findById(request.id) catch {
        try res.json(.{ .@"error" = "not found" }, .{});
        return;
    } orelse return;

    try res.json(Response{
        .id = event.id,
        .name = std.mem.sliceTo(&event.name, 0),
    }, .{});
}

fn handlePutEvent(req: *httpz.Request, res: *httpz.Response) !void {
    res.status = 200;
    const request = (try req.json(PutRequest)).?;
    std.debug.print("got PUT request: {any}\n", .{request});

    if (request.name.len > repository.name_size) {
        res.status = 400;
        try res.json(.{ .@"error" = "name too big" }, .{});
        return;
    }

    try db.insertByName(request.name);
    const event = try db.findByName(request.name) orelse return;
    try res.json(Response{
        .id = event.id,
        .name = std.mem.sliceTo(&event.name, 0),
    }, .{});
}
