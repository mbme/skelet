#!/usr/bin/env node

const SIGINT = 'SIGINT';

var watch = require('node-watch');
var colors = require('colors');
var Proc = require('child_process');

var log = function () {
    var args = Array.prototype.slice.call(arguments, 0);
    args = ['WATCH: '.black].concat(args);

    console.log.apply(console, args);
};

var args = process.argv.slice(2);
var config = {
    pattern: /\.go$/,
    baseDir: args[0] || '.',
    interval: 5*1000
};
config.buildCommand= 'go build ' + config.baseDir;
config.runCommand= config.baseDir + '/skelet';

var child;
var killApp = function () {
    if (child) {
        child.kill(SIGINT);
        child = null;
    }
};
var startApp = function () {
    child = Proc.spawn(config.runCommand);
    child.stdout.on('data', function(chunk) {
        console.log(chunk.toString().green);
    });
    child.stderr.on('data', function(chunk) {
        console.error(chunk.toString().red);
    });

    log('started app: PID', child.pid);
};

var tryToBuild = function () {
    Proc.exec(config.buildCommand, function (error) {
        if (error) {
            log("build failed");
            return;
        }

        killApp();
        startApp();
    });
};

// CHECKS

log("watching", config.baseDir,
    "for", config.pattern, "changes",
    "with interval", config.interval, "ms");

var changed = false;

// watch fs
watch(config.baseDir, function(file) {
    if (config.pattern.test(file)) {
        changed = true;
    }
});

// check if there were changes
setInterval(function () {
    if (changed) {
        log("changes detected, rebuilding");
        tryToBuild();
    }
    changed = false;
}, config.interval);

// handle Ctrl-C
process.on(SIGINT, function () {
    log('got a SIGINT, closing');
    killApp();
    process.exit(0);
});

// run app first time
log('running app');
tryToBuild();
