<script>
    import { goto } from "$app/navigation";

    let username = "";
    let password = "";

    async function handleLogin(event) {
        event.preventDefault();

        try {
            const response = await fetch("/login", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ username, password }),
                credentials: "include", 
            });

            if (response.ok) {
                const data = await response.json();
                goto("/user/profile"); 
            } else {
                console.error("Login failed");
            }
        } catch (error) {
            console.error("Login failed:", error);
        }
    }
</script>

<div class="position-absolute top-50 start-50 translate-middle">
    <form
        class="row g-3 w-100"
        style="max-width: 500px;"
        on:submit|preventDefault={handleLogin}
    >
        <div class="col-12">
            <label for="username" class="form-label">Username</label>
            <input
                type="text"
                class="form-control"
                id="username"
                bind:value={username}
            />
        </div>
        <div class="col-12">
            <label for="inputPassword" class="form-label">Password</label>
            <input
                type="password"
                class="form-control"
                id="inputPassword"
                bind:value={password}
            />
        </div>
        <div class="col-12">
            <div class="form-check">
                <input
                    class="form-check-input"
                    type="checkbox"
                    id="gridCheck"
                />
                <label class="form-check-label" for="gridCheck">
                    Remember me
                </label>
            </div>
        </div>
        <div class="col-12">
            <button type="submit" class="btn btn-primary w-100">Sign in</button>
        </div>
    </form>
</div>
