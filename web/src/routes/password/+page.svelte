<script lang="ts">
  import Error from '$lib/components/Error.svelte';
  import PasswordInput from '$lib/components/PasswordInput.svelte';
  import { superForm } from 'sveltekit-superforms';
  import SaveIcon from '@lucide/svelte/icons/save';
  import UndoIcon from '@lucide/svelte/icons/undo-2';

  export let data;

  const { form, errors, enhance, message, reset } = superForm(data.form, {
    resetForm: true,
  });

  const mustChangePassword = data.must_change_password;

  function handleCancel() {
    reset();
  }
</script>

<div class="flex min-h-screen flex-col items-center justify-center p-6">
  <h1 class="mb-6 text-4xl font-bold">Reset Password</h1>

  {#if mustChangePassword}
    <div class="mb-6 w-full max-w-md rounded-md border border-yellow-400 bg-yellow-100 p-4 text-yellow-800">
      <strong>Important:</strong> You must reset your password before continuing.
    </div>
  {/if}

  <form
    id="resetPasswordForm"
    class="flex w-full max-w-md flex-col gap-4"
    method="POST"
    use:enhance
    novalidate
  >
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

    <div class="mt-4 flex justify-end gap-2">
      <button
        type="button"
        on:click={handleCancel}
        class="flex items-center border border-gray-300 text-gray-700 px-3 py-1 rounded-md hover:bg-gray-100"
      >
        <UndoIcon class="mr-2 h-5 w-5" />
        Cancel
      </button>
      <button type="submit" class="btn flex items-center">
        <SaveIcon class="mr-2 h-5 w-5" />
        Change
      </button>
    </div>

    {#if $message}
      <p class="text-red-600 mt-4 text-sm">{$message}</p>
    {/if}
  </form>
</div>
