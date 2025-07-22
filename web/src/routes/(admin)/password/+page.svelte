<script lang="ts">
  import Error from '$lib/components/Error.svelte';
  import PasswordInput from '$lib/components/PasswordInput.svelte';
  import { superForm } from 'sveltekit-superforms';
  import { fa } from 'zod/v4/locales';

  let { data } = $props();

  let { form, errors, enhance, message } = superForm(data.form, {
    resetForm: false,
  });
</script>

<div class="flex min-h-screen flex-col items-center justify-center p-6">
  <h1 class="mb-6 text-4xl font-bold">Reset Password</h1>

  <form class="mb-4 flex w-full max-w-md flex-col gap-2" method="post" use:enhance novalidate>
    <div class="flex flex-col gap-2">
      <label for="old_password">Current Password</label>
      <PasswordInput
        id="old_password"
        name="old_password"
        bind:value={$form.old_password}
        error={!!$errors.old_password}
      />
      <Error>{$errors.old_password}</Error>
    </div>

    <div class="flex flex-col gap-2">
      <label for="new_password">New Password</label>
      <PasswordInput
        id="new_password"
        name="new_password"
        bind:value={$form.new_password}
        error={!!$errors.new_password}
      />
      <Error>{$errors.new_password}</Error>
    </div>

    <div class="flex flex-col gap-2">
      <label for="confirm_password">Confirm New Password</label>
      <PasswordInput
        id="confirm_password"
        name="confirm_password"
        bind:value={$form.confirm_password}
        error={!!$errors.confirm_password}
      />
      <Error>{$errors.confirm_password}</Error>
    </div>

    <button type="submit" class="btn mt-4">Change Password</button>
  </form>

  {#if $message}
    <p class="text-red-600 mt-2 text-sm">{$message}</p>
  {/if}
</div>
