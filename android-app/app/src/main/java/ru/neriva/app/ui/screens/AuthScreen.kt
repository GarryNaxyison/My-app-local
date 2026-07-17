package ru.neriva.app.ui.screens

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Visibility
import androidx.compose.material.icons.filled.VisibilityOff
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.text.input.PasswordVisualTransformation
import androidx.compose.ui.text.input.VisualTransformation
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import kotlinx.coroutines.launch
import ru.neriva.app.NERIVAApp
import ru.neriva.app.ui.components.AnimatedBackground
import coil.compose.AsyncImage
import coil.request.ImageRequest

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun AuthScreen(onAuthSuccess: () -> Unit) {
    val app = NERIVAApp.instance
    val scope = rememberCoroutineScope()
    val context = LocalContext.current
    var mode by remember { mutableStateOf("login") }
    var login by remember { mutableStateOf("") }
    var password by remember { mutableStateOf("") }
    var passwordConfirm by remember { mutableStateOf("") }
    var referral by remember { mutableStateOf("") }
    var showPassword by remember { mutableStateOf(false) }
    var isLoading by remember { mutableStateOf(false) }
    var error by remember { mutableStateOf<String?>(null) }

    Scaffold(
        containerColor = MaterialTheme.colorScheme.background,
        topBar = {
            CenterAlignedTopAppBar(title = { Text("NERIVA", fontWeight = FontWeight.Bold) })
        }
    ) { padding ->
        AnimatedBackground(modifier = Modifier.padding(padding)) {
            Column(
                modifier = Modifier
                    .fillMaxSize()
                    .verticalScroll(rememberScrollState())
                    .padding(horizontal = 24.dp),
                horizontalAlignment = Alignment.CenterHorizontally,
            ) {
                Spacer(Modifier.height(24.dp))

                // Hero image from server
                AsyncImage(
                    model = ImageRequest.Builder(context)
                        .data("https://neriva.ru/app/assets/auth-login-hero.png")
                        .crossfade(true)
                        .build(),
                    contentDescription = "NERIVA hero",
                    modifier = Modifier.size(120.dp),
                )

                Spacer(Modifier.height(16.dp))
                Text("AI Language Tutor", style = MaterialTheme.typography.headlineSmall, color = MaterialTheme.colorScheme.onBackground)
                Text("Premium daily practice", style = MaterialTheme.typography.bodyMedium,
                    color = MaterialTheme.colorScheme.onSurfaceVariant)
                Spacer(Modifier.height(32.dp))

                // Mode toggle
                Row(modifier = Modifier.fillMaxWidth(), horizontalArrangement = Arrangement.spacedBy(8.dp)) {
                    FilterChip(selected = mode == "login", onClick = { mode = "login" }, label = { Text("Sign In") }, modifier = Modifier.weight(1f))
                    FilterChip(selected = mode == "register", onClick = { mode = "register" }, label = { Text("Register") }, modifier = Modifier.weight(1f))
                }
                Spacer(Modifier.height(20.dp))

                OutlinedTextField(
                    value = login, onValueChange = { login = it },
                    label = { Text("Login") },
                    modifier = Modifier.fillMaxWidth(),
                    singleLine = true,
                )
                Spacer(Modifier.height(12.dp))

                OutlinedTextField(
                    value = password, onValueChange = { password = it },
                    label = { Text("Password") },
                    modifier = Modifier.fillMaxWidth(),
                    singleLine = true,
                    visualTransformation = if (showPassword) VisualTransformation.None else PasswordVisualTransformation(),
                    keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Password),
                    trailingIcon = {
                        IconButton(onClick = { showPassword = !showPassword }) {
                            Icon(if (showPassword) Icons.Default.VisibilityOff else Icons.Default.Visibility, contentDescription = null)
                        }
                    }
                )

                if (mode == "register") {
                    Spacer(Modifier.height(12.dp))
                    OutlinedTextField(
                        value = passwordConfirm, onValueChange = { passwordConfirm = it },
                        label = { Text("Confirm Password") },
                        modifier = Modifier.fillMaxWidth(),
                        singleLine = true,
                        visualTransformation = if (showPassword) VisualTransformation.None else PasswordVisualTransformation(),
                        keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Password),
                    )
                    Spacer(Modifier.height(12.dp))
                    OutlinedTextField(
                        value = referral, onValueChange = { referral = it },
                        label = { Text("Referral code (optional)") },
                        modifier = Modifier.fillMaxWidth(),
                        singleLine = true,
                    )
                }

                Spacer(Modifier.height(8.dp))
                error?.let {
                    Text(it, color = MaterialTheme.colorScheme.error, style = MaterialTheme.typography.bodySmall,
                        textAlign = TextAlign.Center, modifier = Modifier.fillMaxWidth())
                    Spacer(Modifier.height(8.dp))
                }

                Spacer(Modifier.height(16.dp))
                Button(
                    onClick = {
                        isLoading = true; error = null
                        scope.launch {
                            try {
                                if (mode == "login") {
                                    app.authRepo.login(login.trim(), password)
                                } else {
                                    app.authRepo.register(login.trim(), password, passwordConfirm,
                                        referralCode = referral.trim().ifBlank { null })
                                }
                                onAuthSuccess()
                            } catch (e: Exception) {
                                error = e.message ?: "Login failed"
                            } finally { isLoading = false }
                        }
                    },
                    modifier = Modifier.fillMaxWidth().height(52.dp),
                    enabled = !isLoading && login.isNotBlank() && password.isNotBlank(),
                ) {
                    if (isLoading) CircularProgressIndicator(modifier = Modifier.size(20.dp), strokeWidth = 2.dp)
                    else Text(if (mode == "login") "Sign In" else "Create Account")
                }

                Spacer(Modifier.height(12.dp))
                OutlinedButton(
                    onClick = {
                        app.encryptedPrefs.edit().putBoolean("logged_in", true).apply()
                        onAuthSuccess()
                    },
                    modifier = Modifier.fillMaxWidth().height(48.dp),
                ) {
                    Text("Demo Mode (test/test)")
                }

                Spacer(Modifier.height(24.dp))
                Text("Telegram login: @NERIVAapp_bot",
                    style = MaterialTheme.typography.bodySmall,
                    color = MaterialTheme.colorScheme.onSurfaceVariant)
                Spacer(Modifier.height(32.dp))
            }
        }
    }
}
