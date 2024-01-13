//go:build android
// +build android

package android_background_service

/*
#cgo LDFLAGS: -landroid

#include <jni.h>
#include <stdlib.h>
*/

/*
	I have compiled a modified gogio that contains the following:

	<application>
		<service android:name="org.gioui.x.worker_android$WorkerService"></service>
	</application>
	<uses-permission android:name="android.permission.FOREGROUND_SERVICE" />
	<uses-permission android:name="android.permission.POST_NOTIFICATIONS"/>
*/

import "C"
import (
	"gioui.org/app"
	"git.wow.st/gmp/jni"
)

//go:generate javac -source 8 -target 8  -bootclasspath $ANDROID_HOME/platforms/android-33/android.jar -d $TEMP/worker_android/classes *.java
//go:generate jar cf worker_android.jar -C $TEMP/worker_android/classes .

func loadWorkerClass(env jni.Env) (jni.Class, error) {
	return jni.LoadClass(env, jni.ClassLoaderFor(env, jni.Object(app.AppContext())), "org/gioui/x/worker_android")
}

func Start() error {
	err := jni.Do(jni.JVMFor(app.JavaVM()), func(env jni.Env) error {
		class, err := loadWorkerClass(env)
		if err != nil {
			return err
		}

		methodId := jni.GetStaticMethodID(env, class, "startService", "(Landroid/content/Context;)V")
		err = jni.CallStaticVoidMethod(env, class, methodId, jni.Value(app.AppContext()))
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func Stop() error {
	err := jni.Do(jni.JVMFor(app.JavaVM()), func(env jni.Env) error {
		class, err := loadWorkerClass(env)
		if err != nil {
			return err
		}

		methodId := jni.GetStaticMethodID(env, class, "stopService", "(Landroid/content/Context;)V")
		err = jni.CallStaticVoidMethod(env, class, methodId, jni.Value(app.AppContext()))
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func IsForegroundRunning() (bool, error) {
	running := false
	err := jni.Do(jni.JVMFor(app.JavaVM()), func(env jni.Env) error {
		class, err := loadWorkerClass(env)
		if err != nil {
			return err
		}

		fieldId := jni.GetStaticFieldID(env, class, "foregroundRunning", "Z")
		running = jni.GetStaticBooleanField(env, class, fieldId)
		return err
	})

	return running, err
}

func StartForeground() error {
	err := jni.Do(jni.JVMFor(app.JavaVM()), func(env jni.Env) error {
		class, err := loadWorkerClass(env)
		if err != nil {
			return err
		}

		methodId := jni.GetStaticMethodID(env, class, "startForeground", "()V")
		err = jni.CallStaticVoidMethod(env, class, methodId, jni.Value(app.AppContext()))
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func StopForeground() error {
	err := jni.Do(jni.JVMFor(app.JavaVM()), func(env jni.Env) error {
		class, err := loadWorkerClass(env)
		if err != nil {
			return err
		}

		methodId := jni.GetStaticMethodID(env, class, "stopForeground", "()V")
		err = jni.CallStaticVoidMethod(env, class, methodId, jni.Value(app.AppContext()))
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func IsRunning() (bool, error) {
	running := false
	err := jni.Do(jni.JVMFor(app.JavaVM()), func(env jni.Env) error {
		class, err := loadWorkerClass(env)
		if err != nil {
			return err
		}

		fieldId := jni.GetStaticFieldID(env, class, "serviceRunning", "Z")
		running = jni.GetStaticBooleanField(env, class, fieldId)
		return nil
	})

	return running, err
}

func IsAvailable() bool {
	return true
}
