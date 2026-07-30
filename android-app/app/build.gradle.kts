plugins {
    alias(libs.plugins.android.application)
    alias(libs.plugins.kotlin.android)
    alias(libs.plugins.kotlin.compose)
    alias(libs.plugins.kotlin.serialization)
    alias(libs.plugins.ksp)
}

fun releaseSigningValue(name: String): String? =
    providers.gradleProperty(name).orElse(providers.environmentVariable(name)).orNull?.trim()?.takeIf { it.isNotEmpty() }

val releaseStoreFile = releaseSigningValue("NERIVA_KEYSTORE_FILE") ?: "release-key.jks"
val releaseStorePassword = releaseSigningValue("NERIVA_KEYSTORE_PASSWORD")
val releaseKeyAlias = releaseSigningValue("NERIVA_KEY_ALIAS")
val releaseKeyPassword = releaseSigningValue("NERIVA_KEY_PASSWORD")

android {
    namespace = "ru.neriva.app"
    compileSdk = 35

    defaultConfig {
        applicationId = "ru.neriva.app"
        minSdk = 26
        targetSdk = 35
        versionCode = 1
        versionName = "1.0.0"
        vectorDrawables { useSupportLibrary = true }
    }

    signingConfigs {
        create("release") {
            storeFile = rootProject.file(releaseStoreFile)
            storePassword = releaseStorePassword.orEmpty()
            keyAlias = releaseKeyAlias.orEmpty()
            keyPassword = releaseKeyPassword.orEmpty()
        }
    }

    buildTypes {
        debug {
            isMinifyEnabled = false
            applicationIdSuffix = ".debug"
        }
        release {
            isMinifyEnabled = true
            isShrinkResources = true
            proguardFiles(getDefaultProguardFile("proguard-android-optimize.txt"), "proguard-rules.pro")
            signingConfig = signingConfigs.getByName("release")
        }
    }

    compileOptions {
        sourceCompatibility = JavaVersion.VERSION_17
        targetCompatibility = JavaVersion.VERSION_17
    }
    kotlinOptions { jvmTarget = "17" }
    buildFeatures { compose = true; buildConfig = true }
    packaging { resources { excludes += "/META-INF/{AL2.0,LGPL2.1}" } }
}

tasks.matching { it.name == "validateSigningRelease" }.configureEach {
    doFirst {
        check(releaseStorePassword != null && releaseKeyAlias != null && releaseKeyPassword != null) {
            "Release signing requires NERIVA_KEYSTORE_PASSWORD, NERIVA_KEY_ALIAS, and NERIVA_KEY_PASSWORD."
        }
    }
}

dependencies {
    implementation(libs.androidx.core.ktx)
    implementation(libs.androidx.lifecycle.runtime.ktx)
    implementation(libs.androidx.lifecycle.viewmodel.compose)
    implementation(libs.androidx.lifecycle.runtime.compose)
    implementation(libs.androidx.activity.compose)
    implementation(libs.androidx.navigation.compose)
    implementation(platform(libs.compose.bom))
    implementation(libs.compose.ui)
    implementation(libs.compose.ui.graphics)
    implementation(libs.compose.ui.tooling.preview)
    implementation(libs.compose.material3)
    implementation(libs.compose.material.icons.extended)
    debugImplementation(libs.compose.ui.tooling)
    implementation(libs.retrofit)
    implementation(libs.retrofit.kotlinx.serialization)
    implementation(libs.okhttp)
    implementation(libs.okhttp.logging)
    implementation(libs.kotlinx.serialization.json)
    implementation(libs.coil.compose)
    implementation(libs.room.runtime)
    implementation(libs.room.ktx)
    ksp(libs.room.compiler)
    implementation(libs.security.crypto)
    implementation(libs.media3.exoplayer)
    implementation(libs.media3.ui)
    testImplementation(libs.junit)
}
