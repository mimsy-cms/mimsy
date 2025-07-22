<script lang="ts">
  import { superForm } from 'sveltekit-superforms';
  import Input from '$lib/components/Input.svelte';
  import PasswordInput from '$lib/components/PasswordInput.svelte';
  import Error from '$lib/components/Error.svelte';

  export let data;

  const { form, errors, enhance, message, reset } = superForm(data.form, {
    resetForm: true
  });
</script>

<div class="flex flex-col gap-6 p-6"
     style="max-width: calc(100vw - 280px - 48px);">
  <h1 class="text-4xl font-medium">Create user</h1>

  <div class="rounded-md border border-gray-300 bg-white p-6 w-full">
    <form class="flex flex-col gap-4" method="POST" use:enhance>
      <div class="flex flex-col gap-2">
        <label for="email">Email</label>
        <Input
          type="email"
          id="email"
          name="email"
          bind:value={$form.email}
          error={!!$errors.email}
        />
        <Error>{$errors.email}</Error>
      </div>

      <div class="flex flex-col gap-2">
        <label for="password">Password (change will be required)</label>
        <PasswordInput
          id="password"
          name="password"
          bind:value={$form.password}
          error={!!$errors.password}
        />
        <Error>{$errors.password}</Error>
      </div>

      <label class="flex items-center gap-2">
        <input type="checkbox" name="isAdmin" bind:checked={$form.isAdmin} />
        Admin
      </label>

      <button type="submit" class="btn mt-4">Create User</button>
    </form>

    {#if $message}
      <p class="text-red-600 mt-4 text-sm">{$message}</p>
    {/if}
  </div>
</div>
