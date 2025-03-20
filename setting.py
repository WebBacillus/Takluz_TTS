import customtkinter
import yaml
from tkinter import messagebox

class ConfigGUI(customtkinter.CTk):
    def __init__(self):
        super().__init__()

        self.title("Configuration Settings")
        self.geometry("800x650")
        customtkinter.set_appearance_mode("System")
        customtkinter.set_default_color_theme("blue")

        self.config_path = "config.yaml"
        self.config = self.load_config()

        self.sidebar_frame = customtkinter.CTkFrame(self, width=140, corner_radius=0)
        self.sidebar_frame.grid(row=0, column=0, rowspan=4, sticky="nsew")
        self.sidebar_frame.grid_rowconfigure(7, weight=1)

        self.sidebar_buttons = {}
        row_num = 0
        sections_order = ["GENERAL", "OBS", "OPEN_AI", "MICROSOFT", "GOOGLE"]  # Define your custom order here
        for section in sections_order:
            if section in self.config.keys():
                button = customtkinter.CTkButton(self.sidebar_frame, text=section, command=lambda s=section: self.show_section(s))
                button.grid(row=row_num, column=0, padx=20, pady=10, sticky="ew")
                self.sidebar_buttons[section] = button
                row_num += 1
        # Add any remaining sections that were not in the custom order
        for section in self.config.keys():
            if section not in sections_order:
                button = customtkinter.CTkButton(self.sidebar_frame, text=section, command=lambda s=section: self.show_section(s))
                button.grid(row=row_num, column=0, padx=20, pady=10, sticky="ew")
                self.sidebar_buttons[section] = button
                row_num += 1

        self.save_button = customtkinter.CTkButton(self, text="Save Config", command=self.save_config)
        self.save_button.grid(row=4, column=1, padx=20, pady=10, sticky="se")


        self.content_frame = customtkinter.CTkFrame(self)
        self.content_frame.grid(row=0, column=1, sticky="nsew", padx=20, pady=20)
        self.content_frame.grid_columnconfigure(1, weight=1)  # Make entry column expand

        self.entry_widgets = {}

        self.show_section("GENERAL")  # Show GENERAL section by default,  OBS also fine
        self.grid_rowconfigure(0, weight=1)
        self.grid_columnconfigure(1, weight=1)


    def load_config(self):
        try:
            with open(self.config_path, "r") as f:
                return yaml.safe_load(f)
        except FileNotFoundError:
            messagebox.showerror("Error", "config.yaml not found!")
            return {}
        except yaml.YAMLError as e:
            messagebox.showerror("Error", f"Error parsing config.yaml:\n{e}")
            return {}


    def save_config(self):
        for key, entry in self.entry_widgets.items():
            section, option = key.split(".")
            try:
                if option.upper() in ("LIMIT", "TIME_LIMIT"):
                    self.config[section][option] = int(entry.get())  # Force integer
                elif option.upper() in ("SPEED", "VOLUME", "SAMPLE_RATE"):
                    try:
                        self.config[section][option] = float(entry.get())
                    except ValueError:
                        # Handle percentage values like "80%"
                        self.config[section][option] = float(entry.get().replace("%", "")) / 100.0 if "%" in entry.get() else entry.get()

                elif option.upper() == "SPEAKER":
                    self.config[section][option] = int(entry.get())
                elif option.upper() == "SAVE_FILE":
                    self.config[section][option] = entry.get().lower() == "true"
                else:
                    self.config[section][option] = entry.get()
            except ValueError:
                messagebox.showerror("Error", f"Invalid value for {option} in {section}.")
                return
            except KeyError:
                # Key might not exist yet, especially for new sections.  Create it.
                print(f"Warning: Key {key} not found in config.  Creating...")
                if section not in self.config:
                    self.config[section] = {}  # Create the section if it's missing
                self.config[section][option] = entry.get()

        try:
            with open(self.config_path, "w") as f:
                yaml.dump(self.config, f, indent=4)
            messagebox.showinfo("Success", "Configuration saved!")
        except Exception as e:
            messagebox.showerror("Error", f"Failed to save configuration:\n{e}")


    def show_section(self, section_name):
        # Clear existing widgets
        for widget in self.content_frame.winfo_children():
            widget.destroy()
        self.entry_widgets.clear()

        # Highlight the selected button, unhighlight others
        for name, button in self.sidebar_buttons.items():
            if name == section_name:
                button.configure(fg_color="gray")  # Highlight current section
            else:
                button.configure(fg_color=customtkinter.ThemeManager.theme["CTkButton"]["fg_color"])

        row_num = 0
        if section_name in self.config:
            for key, value in self.config[section_name].items():
                label = customtkinter.CTkLabel(self.content_frame, text=f"{key}:")
                label.grid(row=row_num, column=0, padx=10, pady=5, sticky="w")  # Left-align labels

                if key.upper() == 'VOICE':
                    # Handle voice selection with dropdowns, depending on the provider
                    options = []
                    if section_name.upper() == "OPEN_AI":
                         options = ["alloy", "ash", "coral", "echo", "fable", "onyx", "nova", "sage", "shimmer"]
                    elif section_name.upper() == "MICROSOFT":
                        options = [ "th-TH-NiwatNeural", "th-TH-AcharaNeural"]
                    if options:  # Only create a ComboBox if there are options
                        entry = customtkinter.CTkComboBox(self.content_frame, values=options)
                        entry.set(value if value in options else options[0])  # Set to value, or first option if invalid
                    else:
                        entry = customtkinter.CTkEntry(self.content_frame)  # Regular entry if no options
                        entry.insert(0, value)
                elif key.upper() == 'MODEL':
                    options = ["tts-1-hd", "tts-1"]
                    entry = customtkinter.CTkComboBox(self.content_frame, values=options)
                    entry.set(value if value in options else options[0])  # Set to value, or first option if invalid
                elif key.upper() == 'SAVE_FILE':
                    # Use radio buttons for boolean values
                    frame = customtkinter.CTkFrame(self.content_frame)
                    frame.grid(row=row_num, column=1, padx=10, pady=5, sticky="w") # Left-align radio buttons

                    var = customtkinter.StringVar(value=str(value).lower())  # Needs to be string for comparison
                    radio_true = customtkinter.CTkRadioButton(frame, text="True", variable=var, value="true")
                    radio_false = customtkinter.CTkRadioButton(frame, text="False", variable=var, value="false")
                    radio_true.pack(side="left", padx=5)
                    radio_false.pack(side="left", padx=5)
                    self.entry_widgets[f"{section_name}.{key}"] = var  # Store the StringVar
                    row_num += 1
                    continue  # Skip creating a regular entry
                elif key.upper() == 'PLAYER': #add combobox player
                    options = ["OBS", "FFPLAY"]
                    entry = customtkinter.CTkComboBox(self.content_frame, values=options)
                    entry.set(value if value in options else options[0]) # set value
                elif key.upper() == 'NAME':
                    options = [
                        "th-TH-Chirp3-HD-Aoede",
                        "th-TH-Chirp3-HD-Charon",
                        "th-TH-Chirp3-HD-Fenrir",
                        "th-TH-Chirp3-HD-Kore",
                        "th-TH-Chirp3-HD-Leda",
                        "th-TH-Chirp3-HD-Orus",
                        "th-TH-Chirp3-HD-Puck",
                        "th-TH-Chirp3-HD-Zephyr",
                        "th-TH-Neural2-C",
                        "th-TH-Standard-A"
                    ]
                    entry = customtkinter.CTkComboBox(self.content_frame, values=options)
                    entry.set(value if value in options else options[0])  # Set to value, or first option if invalid
                else:
                    # Default: standard text entry
                    entry = customtkinter.CTkEntry(self.content_frame)
                    entry.insert(0, value)

                entry.grid(row=row_num, column=1, padx=10, pady=5, sticky="we")  # Keep entry expanding
                self.entry_widgets[f"{section_name}.{key}"] = entry  # Store widget for later access
                row_num += 1
        else:
            # Handle case where section doesn't exist (yet)
            label = customtkinter.CTkLabel(self.content_frame, text=f"Section '{section_name}' not found in config.")
            label.grid(row=0, column=0, padx=10, pady=5)



if __name__ == "__main__":
    app = ConfigGUI()
    app.mainloop()