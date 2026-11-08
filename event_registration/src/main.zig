const std = @import("std");
const Io = std.Io;
const repository = @import("repository");
const server = @import("server");

pub fn main(init: std.process.Init) !void {
    const io = init.io;
    const allocator = init.gpa;

    // var db = try repository.start();
    // std.debug.print("db: {any}\n", .{db});
    // defer db.deinit();

    var server_future = Io.async(io, server.start, .{ io, allocator });

    // const event = try db.find(57);
    // std.debug.print("found event => id: '{}', name: '{s}'\n", .{event.id, event.name});
    _ = try server_future.await(io);
}
