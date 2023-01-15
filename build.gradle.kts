plugins {
    id("antlr")
}

group = "party.para"
version = "1.0-SNAPSHOT"

repositories {
    mavenCentral()
}

dependencies {
    antlr("org.antlr:antlr4:4.11.1")
}

tasks {
    register<GoBuild>("go") {
        dependsOn("generateGrammarSource")
    }

    "jar" {
        dependsOn("go")
    }

    generateGrammarSource {
        arguments = listOf("-Dlanguage=Go", "-visitor")
    }
}

abstract class GoBuild @Inject() constructor() : DefaultTask() {
    @TaskAction
    fun goBuild() {
        initGeneratedSrcGoMod()
        buildSrc()
    }

    private val logger = org.slf4j.LoggerFactory.getLogger(GoBuild::class.java)

    private val goExecutable: File by lazy {
        val goRoot = System.getenv("GOROOT").also {
            if (it?.isBlank() != false) {
                throw IllegalStateException("GOROOT is empty")
            }
        }
        val goBin = File(File(goRoot), "bin")
        val goExecutable = File(goBin, "go")
        logger.info("go executable path: $goExecutable")
        if (!File("$goExecutable.exe").exists() && !goExecutable.exists()) {
            throw IllegalStateException("go executable not exists")
        }

        goExecutable
    }

    private fun initGeneratedSrcGoMod() {
        val buildPath = File(project.projectDir.absolutePath, "build")
        val genSrcPath = File(buildPath, "generated-src")
        val genAntlrSrcPath = File(genSrcPath, "antlr")
        val genMainAntlrSrcPath = File(genAntlrSrcPath, "main")

        val moduleName = "${project.name}-generated"

        val moduleInfo = File(genMainAntlrSrcPath, "go.mod")
        if (!moduleInfo.exists()) {
            val initCmd = arrayOf("$goExecutable", "mod", "init", moduleName)
            val initTask: Process = Runtime.getRuntime().exec(initCmd, null, genMainAntlrSrcPath)
            initTask.printToLog(fail = {
                logger.error("can not init generated src module")
                throw RuntimeException("can not init generated src module")
            })
        }

        val tidyCmd = arrayOf("$goExecutable", "mod", "tidy")
        val tidyTask: Process = Runtime.getRuntime().exec(tidyCmd, null, genMainAntlrSrcPath)
        tidyTask.printToLog(fail = {
            logger.error("can not tidy generated src module")
            throw RuntimeException("can not tidy generated src module")
        })
    }

    private fun buildSrc() {
        val srcPath = File(project.projectDir.absolutePath, "src")
        val srcMainPath = File(srcPath, "main")
        val goSrcPath = File(srcMainPath, "go")
        logger.info("go code path: $goSrcPath")

        val outputPath = File(project.projectDir.absolutePath, "build")
        val outputName = project.name.let {
            if (System.getProperty("os.name").toLowerCase().contains("windows")) {
                "$it.exe"
            } else {
                it
            }
        }

        val buildCmd = arrayOf("$goExecutable", "build", "-o", "$outputPath/$outputName")
        val task: Process = Runtime.getRuntime().exec(buildCmd, null, goSrcPath)
        task.printToLog(fail = {
            logger.error("go build failed")
            throw RuntimeException("go build failed")
        })
    }

    private fun Process.printToLog(fail: () -> Unit) {
        errorReader(charset("UTF-8")).also {
            try {
                logger.info(it.readLine())
            } catch (_: Throwable) {
            }
        }
        waitFor()
        if (exitValue() != 0) {
            fail()
        }
    }

}
