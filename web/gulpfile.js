'use strict'

const gulp    = require('gulp')
const gutil   = require('gulp-util')
const babel   = require('gulp-babel')
const less    = require('gulp-less')
const compass = require('gulp-compass')
const path    = require('path')
const del     = require('del')

const baseSourcePath = 'src'
const baseDestPath   = 'assets/core'
const sourcePathJS   = `${baseSourcePath}/js/**.js`
const sourcePathSASS = `${baseSourcePath}/sass/**.sass`
const sourcePathLESS = `${baseSourcePath}/less/**.less`
const destPathJS     = `${baseDestPath}/js`
const destPathCSS    = `${baseDestPath}/css`

const noop = (() => {})

gulp.task('clean', () => {
	del([`${baseDestPath}/*/*`])
})

gulp.task('babel', () => {
	gulp.src(sourcePathJS)
		.pipe(babel({ presets: ['es2015'] }).on('error', gutil.log))
		.pipe(gulp.dest(destPathJS))
})

gulp.task('babel:watch', ['babel'], () => {
	gulp.watch(sourcePathJS, ['babel'])
})

gulp.task('less', () => {
	gulp.src(sourcePathLESS)
		.pipe(less({ paths: ['./src/less'] }).on('error', gutil.log))
		.pipe(gulp.dest(destPathCSS))
})

gulp.task('less:watch', ['less'], () => {
	gulp.watch(sourcePathLESS, ['less'])
})

gulp.task('compass', () => {
	gulp.src(sourcePathSASS)
		.pipe(compass({ css: destPathCSS, sass: './src/sass' }).on('error', gutil.log))
		.pipe(gulp.dest(destPathCSS))
})

gulp.task('compass:watch', ['compass'], () => {
	gulp.watch(sourcePathSASS, ['compass'])
})

// let tasks = ['clean', 'babel:watch', 'compass:watch']
let tasks = ['clean', 'babel:watch', 'less:watch']
gulp.task('default', tasks, noop)
