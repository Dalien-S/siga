const std = @import("std");

pub fn build(b: *std.Build) void {
    const target = b.standardTargetOptions(.{});
    const optimize = b.standardOptimizeOption(.{});

    const sqlite = b.dependency("sqlite", .{
        .target = target,
        .optimize = optimize,
    });

    const httpz = b.dependency("httpz", .{
        .target = target,
        .optimize = optimize,
    });

    const exe = b.addExecutable(.{
        .name = "registration",
        .root_module = b.createModule(.{
            .root_source_file = b.path("src/main.zig"),
            .optimize = optimize,
            .target = target,
        }),
    });

    const repo = b.addLibrary(.{
        .name = "repository",
        .root_module = b.createModule(.{
            .root_source_file = b.path("src/repository/repository.zig"),
            .target = target,
            .optimize = optimize,
        }),
    });

    const server = b.addLibrary(.{
        .name = "server",
        .root_module = b.createModule(.{
            .root_source_file = b.path("src/server/server.zig"),
            .target = target,
            .optimize = optimize,
        }),
    });

    exe.root_module.addImport("repository", repo.root_module);
    exe.root_module.addImport("server", server.root_module);
    repo.root_module.addImport("sqlite", sqlite.module("sqlite"));
    server.root_module.addImport("httpz", httpz.module("httpz"));
    server.root_module.addImport("repository", repo.root_module);

    b.installArtifact(exe);

    const run_exe = b.addRunArtifact(exe);
    const run_step = b.step("run", "Run the program");
    run_step.dependOn(&run_exe.step);
}
