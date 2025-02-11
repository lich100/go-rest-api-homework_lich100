package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strconv"
    "github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
    ID           string   `json:"id"`
    Description  string   `json:"description"`
    Note         string   `json:"note"`
    Applications []string `json:"applications"`
}

var tasks = map[string]Task{
    "1": {
        ID:          "1",
        Description: "Сделать финальное задание темы REST API",
        Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
        Applications: []string{
            "VS Code",
            "Terminal",
            "git",
        },
    },
    "2": {
        ID:          "2",
        Description: "Протестировать финальное задание с помощью Postmen",
        Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
        Applications: []string{
            "VS Code",
            "Terminal",
            "git",
            "Postman",
        },
    },
}

// getAllTasksHandler возвращает все задачи
func getAllTasksHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    allTasks := make([]Task, 0, len(tasks))
    for _, task := range tasks {
        allTasks = append(allTasks, task)
    }

    json.NewEncoder(w).Encode(allTasks)
}

// createTaskHandler создаёт новую задачу
func createTaskHandler(w http.ResponseWriter, r *http.Request) {
    var newTask Task
    err := json.NewDecoder(r.Body).Decode(&newTask)
    if err != nil {
        http.Error(w, "Не удалось распарсить тело запроса", http.StatusBadRequest)
        return
    }
    defer r.Body.Close()

    newID := strconv.Itoa(len(tasks) + 1)
    tasks[newID] = newTask
    w.WriteHeader(http.StatusCreated)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(newTask)
}

// getTaskByIDHandler возвращает задачу по заданному ID
func getTaskByIDHandler(w http.ResponseWriter, r *http.Request) {
    taskID := chi.URLParam(r, "id")

    task, ok := tasks[taskID]
    if !ok {
        http.Error(w, "Задача не найдена", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(task)
}

// deleteTaskByIDHandler удаляет задачу по заданному ID
func deleteTaskByIDHandler(w http.ResponseWriter, r *http.Request) {
    taskID := chi.URLParam(r, "id")

    _, ok := tasks[taskID]
    if !ok {
        http.Error(w, "Задача не найдена", http.StatusNotFound)
        return
    }

    delete(tasks, taskID)
    w.WriteHeader(http.StatusOK)
}

func main() {
    r := chi.NewRouter()

    // Регистрация обработчиков
    r.Get("/tasks", getAllTasksHandler)
    r.Post("/tasks", createTaskHandler)
    r.Get("/tasks/{id}", getTaskByIDHandler)
    r.Delete("/tasks/{id}", deleteTaskByIDHandler)

    log.Fatal(http.ListenAndServe(":8080", r))
}