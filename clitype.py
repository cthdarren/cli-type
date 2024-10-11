#!/usr/bin/env python

import os
import configparser
from enum import Enum

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

    
    def getconfig(self, section, key):
        self.config.read(self.configpath)
        print(self.config[section][key])

    def setconfig(self, section, key, value):
        self.config.read(self.configpath)
        self.config[section][key] = value
        with open(self.configpath, "w") as cfgfile:
            self.config.write(cfgfile)

class TypingTest:
    def __init__(self):
        self.config = ConfigurationHandler()
        self.mode = Mode.WORD


ConfigurationHandler().getconfig("Test", "Mode")
ConfigurationHandler().setconfig("Test", "sode", "asdf")
ConfigurationHandler().setconfig("Test", "ss$de", "asdl")
