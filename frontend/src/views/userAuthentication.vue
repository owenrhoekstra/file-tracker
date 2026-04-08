<script setup lang="ts">
import { ref } from 'vue';
import { zodResolver } from '@primevue/forms/resolvers/zod';
import { useToast } from "primevue/usetoast";
import { z } from 'zod';
import router from "../router";

const toast = useToast();
const initialValues = ref({
  username: '',
  email: ''
});

const resolver = ref(zodResolver(
    z.object({
      username: z.string().min(1, { message: 'Username is required.' }),
      email: z.string()
          .min(1, { message: 'Email is required.' })
          .email({ message: 'Invalid email address.' })
    })
));

const onFormSubmit = (event: { valid: boolean }) => {
  const { valid } = event;
  if (valid) {
    toast.add({ severity: 'success', summary: 'Form is submitted.', life: 3000 });
  }
};
</script>

  <template>
    <div class="flex items-center justify-center min-h-screen">
      <Form
          class="flex flex-col gap-4 w-full max-w-xs px-4 mx-auto"
          v-slot="$form"
          :resolver="resolver"
          :initialValues="initialValues"
          @submit="onFormSubmit"
      >
        <div class="flex flex-col gap-1">
          <InputText name="username" type="text" placeholder="Username" fluid />
          <Message v-if="$form.username?.invalid" severity="error" size="small" variant="simple">{{ $form.username.error?.message }}</Message>
        </div>
        <div class="flex flex-col gap-1">
          <InputText name="email" type="text" placeholder="Email" fluid />
          <Message v-if="$form.email?.invalid" severity="error" size="small" variant="simple">{{ $form.email.error?.message }}</Message>
        </div>
        <Button type="submit" severity="secondary" label="Submit" />
        <Button id="forward" label="Print" @click="() => router.push('/print')" />
      </Form>
    </div>
  </template>

<style scoped>

</style>