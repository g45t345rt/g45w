package org.gioui.x;

import android.content.Context;

//import android.util.Log;

import android.app.Notification;
import android.app.NotificationChannel;
import android.app.NotificationManager;
import android.app.Service;
import android.content.pm.ServiceInfo;
import android.content.Intent;
import android.os.Build;
import android.os.Binder;
import android.os.IBinder;
import android.graphics.Bitmap;
import android.graphics.Canvas;
import android.graphics.Color;
import android.graphics.Paint;
import android.graphics.drawable.Icon;
import android.content.ComponentName;
import android.content.ServiceConnection;

public class worker_android {
  public static boolean serviceRunning;
  public static boolean foregroundRunning;
  private static WorkerService workerService;

  public static class WorkerService extends Service {
    private static final String CHANNEL_ID = "ForegroundServiceChannel";

    @Override
    public void onCreate() {
      super.onCreate();
      workerService = this;
    }

    @Override
    public int onStartCommand(Intent intent, int flags, int startId) {
      // startWithForeground();
      serviceRunning = true;
      return START_STICKY;
    }

    public void startWithForeground() {
      createNotificationChannel();
      Notification notification = createNotification();
      startForeground(1, notification);
      foregroundRunning = true;
    }

    public void removeFromForeground() {
      stopForeground(true);
      foregroundRunning = false;
    }

    @Override
    public void onDestroy() {
      super.onDestroy();
      serviceRunning = false;
      foregroundRunning = false;
      workerService = null;
    }

    @Override
    public IBinder onBind(Intent intent) {
      return null;
    }

    private void createNotificationChannel() {
      if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
        NotificationChannel serviceChannel = new NotificationChannel(
            CHANNEL_ID,
            "Foreground Service Channel",
            NotificationManager.IMPORTANCE_DEFAULT);

        NotificationManager manager = getSystemService(NotificationManager.class);
        manager.createNotificationChannel(serviceChannel);
      }
    }

    private Notification createNotification() {
      Icon icon = Icon.createWithBitmap(whiteIcon());
      return new Notification.Builder(this, CHANNEL_ID)
          .setContentTitle("G45W")
          .setContentText("Running in the background.")
          .setSmallIcon(icon)
          .build();
    }

    private Bitmap whiteIcon() {
      Bitmap bitmap = Bitmap.createBitmap(64, 64, Bitmap.Config.ARGB_8888);
      Canvas canvas = new Canvas(bitmap);
      Paint paint = new Paint();
      paint.setColor(Color.WHITE);
      canvas.drawRect(0, 0, 64, 64, paint);
      return bitmap;
    }
  }

  public static void startForeground() {
    if (workerService != null) {
      workerService.startWithForeground();
    }
  }

  public static void stopForeground() {
    if (workerService != null) {
      workerService.removeFromForeground();
    }
  }

  public static void startService(Context context) {
    Intent intent = new Intent(context, WorkerService.class);
    context.startService(intent);
  }

  public static void stopService(Context context) {
    Intent intent = new Intent(context, WorkerService.class);
    context.stopService(intent);
  }
}
