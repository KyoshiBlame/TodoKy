const defaultApiUrl = window.location.protocol.startsWith("http")
  ? `${window.location.origin}/api/v1`
  : "http://localhost:5050/api/v1";

const state = { apiUrl: localStorage.getItem("todoky-api-url") || defaultApiUrl, users: [], tasks: [], selectedUserId: Number(localStorage.getItem("todoky-user-id")) || null, filter: "active" };
const $ = (selector) => document.querySelector(selector);

function apiUrl(path, query = {}) {
  const url = new URL(`${state.apiUrl.replace(/\/$/, "")}${path}`);
  Object.entries(query).forEach(([key, value]) => value !== undefined && value !== null && value !== "" && url.searchParams.set(key, value));
  return url;
}

async function request(path, options = {}, query) {
  const response = await fetch(apiUrl(path, query), { headers: { "Content-Type": "application/json", ...(options.headers || {}) }, ...options });
  if (response.status === 204) return null;
  const data = await response.json().catch(() => null);
  if (!response.ok) throw new Error(data?.message || data?.error || `Ошибка сервера (${response.status})`);
  return data;
}

const api = {
  getUsers: () => request("/users", {}, { limit: 100, offset: 0 }),
  getUser: (id) => request(`/users/${id}`),
  createUser: (data) => request("/users", { method: "POST", body: JSON.stringify(data) }),
  patchUser: (id, data) => request(`/users/${id}`, { method: "PATCH", body: JSON.stringify(data) }),
  deleteUser: (id) => request(`/users/${id}`, { method: "DELETE" }),
  getTasks: (userId) => request("/tasks", {}, { user_id: userId, limit: 100, offset: 0 }),
  getTask: (id) => request(`/tasks/${id}`),
  createTask: (data) => request("/tasks", { method: "POST", body: JSON.stringify(data) }),
  patchTask: (id, data) => request(`/tasks/${id}`, { method: "PATCH", body: JSON.stringify(data) }),
  deleteTask: (id) => request(`/tasks/${id}`, { method: "DELETE" }),
  getStatistics: (userId) => request("/statistics", {}, { user_id: userId }),
};

function showError(error) { $("#status-message").textContent = error.message || "Не удалось выполнить запрос"; }
function clearError() { $("#status-message").textContent = ""; }
function initials(name) { return name.split(/\s+/).map((part) => part[0]).join("").slice(0, 2).toUpperCase(); }
function formatDate(value) { if (!value) return ""; const date = new Date(value); return Number.isNaN(date) ? value : date.toLocaleDateString("ru-RU", { day: "numeric", month: "short" }); }

function renderUsers() {
  const list = $("#users-list");
  list.replaceChildren(...state.users.map((user) => {
    const button = document.createElement("button");
    button.className = `user-item${user.id === state.selectedUserId ? " active" : ""}`;
    button.innerHTML = `<span class="avatar">${initials(user.full_name)}</span><span class="user-name">${user.full_name}</span>`;
    button.onclick = () => selectUser(user.id);
    button.ondblclick = () => openUserDialog(user.id);
    return button;
  }));
}

function renderStatistics(stat = {}) {
  $("#metric-created").textContent = stat.tasks_created ?? "—";
  $("#metric-completed").textContent = stat.tasks_completed ?? "—";
  $("#metric-rate").textContent = stat.tasks_completed_rate == null ? "—" : `${Math.round(stat.tasks_completed_rate)}%`;
  $("#metric-average").textContent = stat.tasks_avg_completion_time || "—";
}

function renderTasks() {
  const active = state.tasks.filter((task) => !task.completed);
  const completed = state.tasks.filter((task) => task.completed);
  $("#active-count").textContent = active.length;
  $("#completed-count").textContent = completed.length;
  const shown = state.filter === "active" ? active : state.filter === "completed" ? completed : state.tasks;
  const list = $("#tasks-list");
  list.replaceChildren(...shown.map((task) => {
    const item = document.createElement("article");
    item.className = `task${task.completed ? " completed" : ""}`;
    item.innerHTML = `<input class="task-toggle" type="checkbox" aria-label="Изменить статус задачи" ${task.completed ? "checked" : ""}><div class="task-content"><h3></h3><p></p><div class="task-meta"></div></div><div class="task-actions"><button type="button" data-action="edit">Изменить</button><button type="button" class="delete" data-action="delete">Удалить</button></div>`;
    item.querySelector("h3").textContent = task.title;
    const description = item.querySelector("p"); description.textContent = task.description || ""; description.hidden = !task.description;
    item.querySelector(".task-meta").textContent = task.created_at ? `Создано ${formatDate(task.created_at)}` : "";
    item.querySelector(".task-toggle").onchange = () => updateTask(task.id, { completed: !task.completed });
    item.querySelector('[data-action="edit"]').onclick = () => openTaskDialog(task.id);
    item.querySelector('[data-action="delete"]').onclick = () => deleteTask(task.id);
    return item;
  }));
  $("#empty-state").hidden = shown.length !== 0 || !state.selectedUserId;
}

async function selectUser(id) {
  state.selectedUserId = id; localStorage.setItem("todoky-user-id", id); renderUsers(); await refreshDashboard();
}

async function refreshDashboard() {
  if (!state.selectedUserId) { state.tasks = []; renderTasks(); renderStatistics(); $("#add-task-button").disabled = true; $("#page-title").textContent = "Выберите участника"; return; }
  clearError(); $("#add-task-button").disabled = false;
  const user = state.users.find((item) => item.id === state.selectedUserId);
  $("#page-title").textContent = user ? `Задачи: ${user.full_name}` : "Мои задачи";
  try { const [tasks, statistics] = await Promise.all([api.getTasks(state.selectedUserId), api.getStatistics(state.selectedUserId)]); state.tasks = tasks || []; renderTasks(); renderStatistics(statistics); } catch (error) { showError(error); }
}

async function loadUsers() {
  clearError();
  try { state.users = await api.getUsers() || []; if (state.selectedUserId && !state.users.some((user) => user.id === state.selectedUserId)) state.selectedUserId = null; if (!state.selectedUserId && state.users[0]) state.selectedUserId = state.users[0].id; renderUsers(); await refreshDashboard(); } catch (error) { showError(error); }
}

function openTaskDialog(id) {
  const dialog = $("#task-dialog"); const form = $("#task-form"); form.reset(); $("#task-id").value = "";
  if (id) { api.getTask(id).then((task) => { $("#task-dialog-title").textContent = "Изменить задачу"; $("#task-id").value = task.id; $("#task-title").value = task.title; $("#task-description").value = task.description || ""; dialog.showModal(); }).catch(showError); return; }
  $("#task-dialog-title").textContent = "Новая задача"; dialog.showModal();
}

async function updateTask(id, data) { try { await api.patchTask(id, data); await refreshDashboard(); } catch (error) { showError(error); } }
async function deleteTask(id) { if (!confirm("Удалить эту задачу?")) return; try { await api.deleteTask(id); await refreshDashboard(); } catch (error) { showError(error); } }

function openUserDialog(id) {
  const dialog = $("#user-dialog"); const form = $("#user-form"); form.reset(); $("#user-id").value = "";
  $("#delete-user-button").hidden = !id;
  if (id) { api.getUser(id).then((user) => { $("#user-dialog-title").textContent = "Изменить пользователя"; $("#user-id").value = user.id; $("#user-name").value = user.full_name; $("#user-phone").value = user.phone_number || ""; dialog.showModal(); }).catch(showError); return; }
  $("#user-dialog-title").textContent = "Новый пользователь"; dialog.showModal();
}

$("#task-form").addEventListener("submit", async (event) => { event.preventDefault(); const id = Number($("#task-id").value); const data = { title: $("#task-title").value.trim(), description: $("#task-description").value.trim() || null }; try { if (id) await api.patchTask(id, data); else await api.createTask({ ...data, author_user_id: state.selectedUserId }); $("#task-dialog").close(); await refreshDashboard(); } catch (error) { showError(error); } });
$("#user-form").addEventListener("submit", async (event) => { event.preventDefault(); const id = Number($("#user-id").value); const data = { full_name: $("#user-name").value.trim(), phone_number: $("#user-phone").value.trim() || null }; try { const user = id ? await api.patchUser(id, data) : await api.createUser(data); $("#user-dialog").close(); state.selectedUserId = user.id; await loadUsers(); } catch (error) { showError(error); } });
$("#api-form").addEventListener("submit", (event) => { event.preventDefault(); state.apiUrl = $("#api-url").value.trim().replace(/\/$/, ""); localStorage.setItem("todoky-api-url", state.apiUrl); $("#api-dialog").close(); loadUsers(); });
$("#delete-user-button").onclick = async () => { const id = Number($("#user-id").value); if (!id || !confirm("Удалить пользователя?")) return; try { await api.deleteUser(id); $("#user-dialog").close(); state.selectedUserId = null; await loadUsers(); } catch (error) { showError(error); } };
$("#add-user-button").onclick = () => openUserDialog(); $("#add-task-button").onclick = () => openTaskDialog(); $("#refresh-button").onclick = refreshDashboard;
$("#api-settings-button").onclick = () => { $("#api-url").value = state.apiUrl; $("#api-dialog").showModal(); };
document.querySelectorAll("[data-close-dialog]").forEach((button) => { button.onclick = () => $(`#${button.dataset.closeDialog}`).close(); });
document.querySelectorAll(".tab").forEach((button) => { button.onclick = () => { state.filter = button.dataset.filter; document.querySelectorAll(".tab").forEach((tab) => { const active = tab === button; tab.classList.toggle("active", active); tab.setAttribute("aria-selected", active); }); renderTasks(); }; });
$("#current-date").textContent = new Intl.DateTimeFormat("ru-RU", { weekday: "long", day: "numeric", month: "long" }).format(new Date());
loadUsers();
