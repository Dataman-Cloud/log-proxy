var path = require('path');
var gulp = require('gulp');
var conf = require('./conf');

var $ = require('gulp-load-plugins')({
    pattern: ['gulp-*', 'main-bower-files', 'del']
});

gulp.task('copy', function () {
    var fonts = ['bower_components/font-awesome/fonts/*'];

    gulp.src(fonts)
        .pipe(gulp.dest('dist/fonts'));
    gulp.src('pics/**')
        .pipe(gulp.dest('dist/pics'));
});

//app html to js
gulp.task('template-min-app', function () {
    return gulp.src('src/app.html')
        .pipe($.minifyHtml({
            empty: true,
            spare: true,
            quotes: true
        }))
        .pipe($.angularTemplatecache('templateCacheHtmlApp.js', {
            module: 'app',
            root: '/src/'
        }))
        .pipe(gulp.dest('dist/src/'));
});

//directives html to js
gulp.task('template-min-directives', ['template-min-app'], function () {
    return gulp.src('src/directives/**/*.html')
        .pipe($.minifyHtml({
            empty: true,
            spare: true,
            quotes: true
        }))
        .pipe($.angularTemplatecache('templateCacheHtmlDirectives.js', {
            module: 'app',
            root: '/src/directives'
        }))
        .pipe(gulp.dest('dist/src/'));
});

//utils html to js
gulp.task('template-min-utils', ['template-min-directives'], function () {
    return gulp.src('src/utils/**/*.html')
        .pipe($.minifyHtml({
            empty: true,
            spare: true,
            quotes: true
        }))
        .pipe($.angularTemplatecache('templateCacheHtmlUtils.js', {
            module: 'app',
            root: '/src/utils'
        }))
        .pipe(gulp.dest('dist/src/'));
});

//dashboard html to js
gulp.task('template-min-dashboard', ['template-min-utils'], function () {
    return gulp.src('src/dashboard/**/*.html')
        .pipe($.minifyHtml({
            empty: true,
            spare: true,
            quotes: true
        }))
        .pipe($.angularTemplatecache('templateCacheHtmlDashboard.js', {
            module: 'app',
            root: '/src/dashboard'
        }))
        .pipe(gulp.dest('dist/src/'));
});

//monitor html to js
gulp.task('template-min-monitor', ['template-min-dashboard'], function () {
    return gulp.src('src/monitor/**/*.html')
        .pipe($.minifyHtml({
            empty: true,
            spare: true,
            quotes: true
        }))
        .pipe($.angularTemplatecache('templateCacheHtmlMonitor.js', {
            module: 'app',
            root: '/src/monitor'
        }))
        .pipe(gulp.dest('dist/src/'));
});

//log html to js
gulp.task('template-min-log', ['template-min-monitor'], function () {
    return gulp.src('src/log/**/*.html')
        .pipe($.minifyHtml({
            empty: true,
            spare: true,
            quotes: true
        }))
        .pipe($.angularTemplatecache('templateCacheHtmlLog.js', {
            module: 'app',
            root: '/src/log'
        }))
        .pipe(gulp.dest('dist/src/'));
});

//alert html to js
gulp.task('template-min-alert', ['template-min-log'], function () {
    return gulp.src('src/alert/**/*.html')
        .pipe($.minifyHtml({
            empty: true,
            spare: true,
            quotes: true
        }))
        .pipe($.angularTemplatecache('templateCacheHtmlAlert.js', {
            module: 'app',
            root: '/src/alert'
        }))
        .pipe(gulp.dest('dist/src/'));
});

gulp.task('ng-annotate', ['template-min-alert'], function () {
    return gulp.src('src/**/*.js')
        .pipe($.ngAnnotate({add: true}))
        .pipe(gulp.dest('dist/src'))
});

gulp.task('html-replace', ['ng-annotate'], function () {
    var templateInjectFile = gulp.src('dist/src/templateCacheHtml*.js', {read: false});
    var templateInjectOptions = {
        starttag: '<!-- inject:template.js  -->',
        addRootSlash: false
    };

    var revAll = new $.revAll();
    return gulp.src('index.html')
        .pipe($.inject(templateInjectFile, templateInjectOptions))
        .pipe($.useref()).on('error', $.util.log)
        .pipe($.if('*.js', $.uglify()))
        .pipe($.if('*.css', $.minifyCss()))
        .pipe(revAll.revision().on('error', $.util.log))
        .pipe($.revHash())
        .pipe(gulp.dest('dist/'))
        .pipe(revAll.manifestFile())
        .pipe(gulp.dest('dist/'));
});

gulp.task('html-rename', ['html-replace'], function () {
    gulp.src('dist/index.*.html')
        .pipe($.rename('index.html').on('error', $.util.log))
        .pipe(gulp.dest('dist/'));
});

gulp.task('clean', ['html-rename'], function () {
    return $.del([path.join(conf.paths.dist, 'index.*.html')]);
});

gulp.task('delete', function () {
    return $.del([path.join(conf.paths.dist, '/'), path.join(conf.paths.tmp, '/')]);
});

gulp.task('build', ['copy', 'clean']);