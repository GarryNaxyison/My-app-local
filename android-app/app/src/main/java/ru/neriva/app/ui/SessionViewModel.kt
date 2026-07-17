package ru.neriva.app.ui

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import ru.neriva.app.NERIVAApp
import ru.neriva.app.data.model.SessionData

data class SessionUiState(
    val isLoading: Boolean = true,
    val isLoggedIn: Boolean? = null,
    val session: SessionData? = null,
    val error: String? = null,
)

class SessionViewModel : ViewModel() {

    private val app get() = NERIVAApp.instance
    private val _state = MutableStateFlow(SessionUiState())
    val state: StateFlow<SessionUiState> = _state.asStateFlow()

    init { checkSession() }

    fun checkSession() {
        viewModelScope.launch {
            _state.value = _state.value.copy(isLoading = true, error = null)
            try {
                val session = app.sessionRepo.getSession()
                _state.value = SessionUiState(isLoading = false, isLoggedIn = session.authenticated, session = session)
            } catch (e: Exception) {
                _state.value = SessionUiState(isLoading = false, isLoggedIn = false, error = e.message)
            }
        }
    }

    fun refresh() = checkSession()

    fun logout() {
        viewModelScope.launch {
            app.authRepo.logout()
            _state.value = SessionUiState(isLoading = false, isLoggedIn = false)
        }
    }
}
