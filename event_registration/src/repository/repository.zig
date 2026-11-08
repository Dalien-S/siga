const std = @import("std");
const sqlite = @import("sqlite");

pub const Repository = @This();

pub const name_size: usize = 128;
pub const Event = struct {
    id: i32,
    name: [name_size:0]u8,
};

db: sqlite.Db,

pub fn init() !Repository {
    std.debug.print("Hello, World!\n", .{});
    var db = try sqlite.Db.init(.{
        .mode = sqlite.Db.Mode{ .File = "events.db" },
        .open_flags = .{ .create = true, .write = true },
        .threading_mode = .MultiThread,
    });

    try db.exec(
        "create table if not exists events(id integer primary key, name text)",
        .{},
        .{},
    );

    const query = "insert or ignore into events(id, name) values(?, ?)";
    var stmt = try db.prepare(query);
    defer stmt.deinit();

    try stmt.exec(.{}, .{
        .id = 57,
        .name = "Main event",
    });

    return Repository{ .db = db };
}

pub fn deinit(r: *Repository) void {
    r.db.deinit();
}

pub fn insertByName(r: *Repository, name: []const u8) !void {
    const query = "insert into events(name) values(?)";
    var stmt = try r.db.prepare(query);
    defer stmt.deinit();

    try stmt.exec(.{}, .{
        .name = name,
        // .name = std.mem.sliceTo(&event.name, 0),
    });
}

pub fn findById(r: *Repository, event_id: i32) !?Event {
    const query = "select id, name from events where id = ?";
    const row = try r.db.one(Event, query, .{}, .{ .id = event_id });

    if (row) |event| {
        std.debug.print("found event by id: {any}\n", .{event});
    }

    return row;
}

pub fn findByName(r: *Repository, event_name: []const u8) !?Event {
    const query = "select id, name from events where name = ?";
    const row = try r.db.one(Event, query, .{}, .{ .name = event_name });

    if (row) |event| {
        std.debug.print("found event by name: {any}\n", .{event});
    }

    return row;
}
