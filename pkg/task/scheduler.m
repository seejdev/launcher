//go:build darwin
// +build darwin

#import <AppKit/AppKit.h>
#import <Foundation/Foundation.h>
#import <Foundation/NSBackgroundActivityScheduler.h>

// Go callbacks
extern void performTask(char*);

void schedule(char* cIdentifier, int repeats, uint64_t interval, void* pActivity) {
  @autoreleasepool {
    NSString* identifier = [NSString stringWithUTF8String:cIdentifier];
    NSBackgroundActivityScheduler* activity =
        [[NSBackgroundActivityScheduler alloc] initWithIdentifier:identifier];

    NSQualityOfService qos = NSQualityOfServiceBackground;
    if (interval < 15) {
      qos = NSQualityOfServiceUserInteractive;
    }
    else if (interval < 60) {
      qos = NSQualityOfServiceUserInitiated;
    }

    activity.repeats = repeats ? YES: NO;
    activity.interval = interval;
    activity.qualityOfService = qos;

    [activity
        scheduleWithBlock:^(NSBackgroundActivityCompletionHandler completion) {
          performTask(cIdentifier);
          completion(NSBackgroundActivityResultFinished);
        }];

    pActivity = activity;
  }
}

void reset(void* p) {
}

void stop(void* p) {
  NSBackgroundActivityScheduler* activity = (NSBackgroundActivityScheduler*)p;
  [activity invalidate];
}