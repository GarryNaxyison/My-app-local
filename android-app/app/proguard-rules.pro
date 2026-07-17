# Retrofit
-keepattributes Signature, InnerClasses, EnclosingMethod
-keepattributes RuntimeVisibleAnnotations
-keepclassmembers interface * { @retrofit2.http.* <methods>; }
# kotlinx.serialization
-keep,includedescriptorclasses class ru.neriva.app.data.model.**$$serializer { *; }
-keepclassmembers class ru.neriva.app.data.model.** { *** Companion; }
# Room
-keep class * extends androidx.room.RoomDatabase
-keep @androidx.room.Entity class *
# OkHttp
-dontwarn okhttp3.**
-dontwarn okio.**
