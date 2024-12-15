const apiUrl = "http://localhost:8080"; // Update with your backend URL

// Add Task
document.getElementById("addTaskForm").addEventListener("submit", async (e) => {
    e.preventDefault();
    const taskName = document.getElementById("newTask").value;
    const response = await fetch(`${apiUrl}/add-task`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({ task_name: taskName, status: false })
    });
    const result = await response.text();
    alert(result);
    document.getElementById("addTaskForm").reset();
});

// Load Tasks
document.getElementById("loadTasksButton").addEventListener("click", async () => {
    const response = await fetch(`${apiUrl}/view-tasks`);
    const tasks = await response.json();
    const tbody = document.querySelector("#tasksTable tbody");
    tbody.innerHTML = ""; // Clear existing tasks
    tasks.forEach((task) => {
        const row = document.createElement("tr");
        row.innerHTML = `
            <td>${task.id}</td>
            <td>${task.task_name}</td>
            <td>${task.status ? "Complete" : "Incomplete"}</td>
            <td>
                <button onclick="markComplete(${task.id}, true)">Complete</button>
                <button onclick="markComplete(${task.id}, false)">Incomplete</button>
                <button onclick="deleteTask(${task.id})">Delete</button>
            </td>
        `;
        tbody.appendChild(row);
    });
});

// Update Task Status
document.getElementById("updateTaskForm").addEventListener("submit", async (e) => {
    e.preventDefault();
    const id = document.getElementById("updateTaskId").value;
    const status = document.getElementById("updateTaskStatus").value === "true";
    const response = await fetch(`${apiUrl}/update-task`, {
        method: "PUT",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({ id: parseInt(id), status })
    });
    const result = await response.text();
    alert(result);
    document.getElementById("updateTaskForm").reset();
});

// Delete Task
document.getElementById("deleteTaskForm").addEventListener("submit", async (e) => {
    e.preventDefault();
    const id = document.getElementById("deleteTaskId").value;
    const response = await fetch(`${apiUrl}/delete-task?id=${id}`, {
        method: "DELETE"
    });
    const result = await response.text();
    alert(result);
    document.getElementById("deleteTaskForm").reset();
});

// Helper functions
async function markComplete(id, status) {
    const response = await fetch(`${apiUrl}/update-task`, {
        method: "PUT",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({ id, status })
    });
    const result = await response.text();
    alert(result);
    document.getElementById("loadTasksButton").click(); // Refresh tasks
}

async function deleteTask(id) {
    const response = await fetch(`${apiUrl}/delete-task?id=${id}`, {
        method: "DELETE"
    });
    const result = await response.text();
    alert(result);
    document.getElementById("loadTasksButton").click(); // Refresh tasks
}
