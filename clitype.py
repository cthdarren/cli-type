#!/usr/bin/env python

import os
import configparser
from enum import Enum

WELCOME_TEXT = """
Welcome to CLI Type. An open source typing test at your fingertips in your cli!.

Press ESC at any time to open the menu.
"""

MENU_TEXT = """
-------------- Menu -------------- 

[  1  ] Change Mode

[  2  ] Change Time/Words

[  3  ] Change Wordlist

[  4  ] View Stats

[  q  ] Quit CLI-type

[ ESC ] Back

----------------------------------
"""

class Mode(Enum):
    WORD = 1
    TIME = 2
    # QUOTE = 3
    # ZEN = 4

class ConfigurationHandler:
    def __init__(self):
        userConfigPath = os.path.expanduser("~") + '.clityperc'
        self.config = configparser.ConfigParser()

        if os.path.exists(userConfigPath):
            self.configpath = userConfigPath
        else:
            if not os.path.exists(".clityperc"):
                open(".clityperc", "a").close()
            self.configpath = ".clityperc"

    
    def getconfig(self, section, key, default_value):
        self.config.read(self.configpath)
        try:
            return self.config[section][key]
        except KeyError:
            print(f"ERROR: Failed to read configuration file at [{section}]{key}, using default values...")
            return default_value

    def setconfig(self, section, key, value):
        self.config.read(self.configpath)
        self.config[section][key] = value
        with open(self.configpath, "w") as cfgfile:
            try: 
                self.config.write(cfgfile)
                return True
            except:
                return False


class TypingTest:
    def __init__(self):
        config = ConfigurationHandler()
        self.mode = int(config.getconfig("Test", "Mode", 1))
        self.time = int(config.getconfig("Test", "Time", 30))
        self.words = int(config.getconfig("Test", "Words", 50))
        self.wordset = int(config.getconfig("Test", "Words", 50))

    def start(self):
        print(WELCOME_TEXT)
        match self.mode:
            case Mode.WORD.value:
                self.word_test()
            case Mode.TIME:
                self.time_test()
            case _:
                print(self.mode)

    def word_test(self):
        return

    def time_test(self):
        return


TypingTest().start()

