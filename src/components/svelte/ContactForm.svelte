<script lang="ts">
  interface Strings {
    name: string;
    email: string;
    message: string;
    namePh: string;
    emailPh: string;
    msgPh: string;
    required: string;
    badEmail: string;
    minChars: string;
    send: string;
    sending: string;
    sent: string;
  }

  interface Props {
    strings?: Strings;
  }

  const { strings = {
    name: 'Name',
    email: 'Email',
    message: 'Message',
    namePh: 'Your name',
    emailPh: 'you@example.com',
    msgPh: 'Tell me about your project…',
    required: 'Required.',
    badEmail: 'Invalid email.',
    minChars: 'At least 10 characters.',
    send: 'Send message →',
    sending: 'Sending…',
    sent: "Message received. I'll be in touch.",
  }}: Props = $props();

  type Field = { value: string; error: string };

  let name: Field = $state({ value: '', error: '' });
  let email: Field = $state({ value: '', error: '' });
  let message: Field = $state({ value: '', error: '' });
  let status: 'idle' | 'sending' | 'sent' = $state('idle');

  function validateEmail(v: string) {
    return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(v);
  }

  function validate(): boolean {
    let ok = true;
    name.error = name.value.trim() ? '' : strings.required;
    if (name.error) ok = false;

    if (!email.value.trim()) { email.error = strings.required; ok = false; }
    else if (!validateEmail(email.value)) { email.error = strings.badEmail; ok = false; }
    else email.error = '';

    message.error = message.value.trim().length >= 10 ? '' : strings.minChars;
    if (message.error) ok = false;

    return ok;
  }

  async function submit(e: SubmitEvent) {
    e.preventDefault();
    if (!validate()) return;
    status = 'sending';
    await new Promise((r) => setTimeout(r, 1000));
    status = 'sent';
  }
</script>

{#if status === 'sent'}
  <p class="sent-msg">{strings.sent}</p>
{:else}
  <form class="form" onsubmit={submit} novalidate>
    <div class="form__row">
      <div class="form__field" class:err={!!name.error}>
        <label for="cf-name" class="form__label label">{strings.name}</label>
        <input id="cf-name" type="text" bind:value={name.value} placeholder={strings.namePh} />
        {#if name.error}<span class="form__err" role="alert">{name.error}</span>{/if}
      </div>
      <div class="form__field" class:err={!!email.error}>
        <label for="cf-email" class="form__label label">{strings.email}</label>
        <input id="cf-email" type="email" bind:value={email.value} placeholder={strings.emailPh} />
        {#if email.error}<span class="form__err" role="alert">{email.error}</span>{/if}
      </div>
    </div>
    <div class="form__field" class:err={!!message.error}>
      <label for="cf-msg" class="form__label label">{strings.message}</label>
      <textarea id="cf-msg" rows="5" bind:value={message.value} placeholder={strings.msgPh}></textarea>
      {#if message.error}<span class="form__err" role="alert">{message.error}</span>{/if}
    </div>
    <button type="submit" class="form__submit label" disabled={status === 'sending'}>
      {status === 'sending' ? strings.sending : strings.send}
    </button>
  </form>
{/if}

<style>
  .form {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }

  .form__row {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1.5rem;
  }

  @media (max-width: 640px) {
    .form__row { grid-template-columns: 1fr; }
  }

  .form__field {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .form__label { color: var(--clr-muted); }

  input, textarea {
    background: transparent;
    border: none;
    border-bottom: 1px solid var(--clr-rule);
    padding: 0.5rem 0;
    font-family: var(--font-body);
    font-size: 0.95rem;
    color: var(--clr-ink);
    transition: border-color var(--transition);
    width: 100%;
    resize: none;
  }

  input::placeholder, textarea::placeholder { color: var(--clr-muted); }

  input:focus, textarea:focus {
    outline: none;
    border-bottom-color: var(--clr-ink);
  }

  .err input, .err textarea { border-bottom-color: var(--clr-error); }

  .form__err {
    font-size: 0.75rem;
    color: var(--clr-error);
  }

  .form__submit {
    align-self: flex-start;
    background: none;
    border: none;
    cursor: pointer;
    color: var(--clr-ink);
    padding: 0.5rem 0;
    border-bottom: 1px solid var(--clr-ink);
    transition: opacity var(--transition);
    font-family: var(--font-body);
  }

  .form__submit:hover { opacity: 0.5; }
  .form__submit:disabled { opacity: 0.4; cursor: not-allowed; }

  .sent-msg {
    font-family: var(--font-display);
    font-style: italic;
    font-size: clamp(1.25rem, 3vw, 2rem);
    color: var(--clr-ink);
  }
</style>
